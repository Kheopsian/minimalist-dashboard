/*!
 * Chart.js Minimal - Custom build for Network Chart
 * Based on Chart.js v4.4.0
 * Contains only: LineController, LinearScale, basic animations, and core functionality
 */
!(function (global, factory) {
  typeof exports === 'object' && typeof module !== 'undefined'
    ? module.exports = factory()
    : typeof define === 'function' && define.amd
    ? define(factory)
    : ((global = typeof globalThis !== 'undefined' ? globalThis : global || self).Chart = factory());
})(this, function () {
  'use strict';

  // Core utilities (minimal subset)
  const PI = Math.PI;
  const TAU = 2 * PI;
  const HALF_PI = PI / 2;
  const RAD_PER_DEG = PI / 180;

  function isArray(value) {
    if (Array.isArray && Array.isArray(value)) return true;
    const type = Object.prototype.toString.call(value);
    return type.slice(0, 7) === '[object' && type.slice(-6) === 'Array]';
  }

  function isObject(value) {
    return value !== null && Object.prototype.toString.call(value) === '[object Object]';
  }

  function isFiniteNumber(value) {
    return (typeof value === 'number' || value instanceof Number) && isFinite(+value);
  }

  function merge(target, source) {
    if (!isObject(target)) return target;
    const sources = isArray(source) ? source : [source];
    
    sources.forEach(source => {
      if (!isObject(source)) return;
      Object.keys(source).forEach(key => {
        const targetValue = target[key];
        const sourceValue = source[key];
        
        if (isObject(targetValue) && isObject(sourceValue)) {
          target[key] = merge(targetValue, sourceValue);
        } else {
          target[key] = sourceValue;
        }
      });
    });
    
    return target;
  }

  // Basic Element class
  class Element {
    constructor() {
      this.x = undefined;
      this.y = undefined;
      this.active = false;
      this.options = undefined;
    }

    tooltipPosition() {
      const {x, y} = this;
      return {x, y};
    }

    hasValue() {
      return isFiniteNumber(this.x) && isFiniteNumber(this.y);
    }
  }

  // Point Element (minimal)
  class PointElement extends Element {
    static id = 'point';
    static defaults = {
      borderWidth: 1,
      hitRadius: 1,
      hoverBorderWidth: 1,
      hoverRadius: 4,
      pointStyle: 'circle',
      radius: 3,
      rotation: 0
    };

    inRange(mouseX, mouseY) {
      const options = this.options;
      const {x, y} = this;
      return Math.pow(mouseX - x, 2) + Math.pow(mouseY - y, 2) < Math.pow(options.hitRadius + options.radius, 2);
    }

    getCenterPoint() {
      const {x, y} = this;
      return {x, y};
    }

    draw(ctx) {
      const options = this.options;
      if (this.skip || options.radius < 0.1) return;
      
      ctx.strokeStyle = options.borderColor;
      ctx.lineWidth = options.borderWidth;
      ctx.fillStyle = options.backgroundColor;
      
      ctx.beginPath();
      ctx.arc(this.x, this.y, options.radius, 0, TAU);
      ctx.fill();
      if (options.borderWidth > 0) {
        ctx.stroke();
      }
    }
  }

  // Line Element (minimal)
  class LineElement extends Element {
    static id = 'line';
    static defaults = {
      borderCapStyle: 'butt',
      borderDash: [],
      borderDashOffset: 0,
      borderJoinStyle: 'miter',
      borderWidth: 3,
      capBezierPoints: true,
      cubicInterpolationMode: 'default',
      fill: false,
      spanGaps: false,
      stepped: false,
      tension: 0
    };

    constructor(cfg) {
      super();
      this.animated = true;
      this.options = undefined;
      this._chart = undefined;
      this._points = undefined;
      this._segments = undefined;
      if (cfg) {
        Object.assign(this, cfg);
      }
    }

    draw(ctx, chartArea) {
      const options = this.options || {};
      const points = this.points || [];
      
      if (points.length < 2) return;
      
      ctx.save();
      ctx.strokeStyle = options.borderColor;
      ctx.lineWidth = options.borderWidth;
      ctx.setLineDash(options.borderDash || []);
      ctx.lineDashOffset = options.borderDashOffset || 0;
      
      ctx.beginPath();
      ctx.moveTo(points[0].x, points[0].y);
      
      for (let i = 1; i < points.length; i++) {
        const point = points[i];
        if (!point.skip) {
          if (options.tension > 0) {
            // Simple curve approximation
            const prev = points[i - 1];
            const cp1x = prev.x + (point.x - prev.x) * 0.5;
            const cp1y = prev.y;
            const cp2x = point.x - (point.x - prev.x) * 0.5;
            const cp2y = point.y;
            ctx.bezierCurveTo(cp1x, cp1y, cp2x, cp2y, point.x, point.y);
          } else {
            ctx.lineTo(point.x, point.y);
          }
        }
      }
      
      if (options.fill) {
        ctx.lineTo(points[points.length - 1].x, chartArea.bottom);
        ctx.lineTo(points[0].x, chartArea.bottom);
        ctx.closePath();
        ctx.fillStyle = options.backgroundColor;
        ctx.fill();
      }
      
      ctx.stroke();
      ctx.restore();
    }

    set points(points) {
      this._points = points;
    }

    get points() {
      return this._points;
    }
  }

  // Basic Scale class
  class Scale extends Element {
    constructor(cfg) {
      super();
      this.id = cfg.id;
      this.type = cfg.type;
      this.options = undefined;
      this.ctx = cfg.ctx;
      this.chart = cfg.chart;
      this.top = undefined;
      this.bottom = undefined;
      this.left = undefined;
      this.right = undefined;
      this.width = undefined;
      this.height = undefined;
      this.min = undefined;
      this.max = undefined;
      this.ticks = [];
      this._startPixel = undefined;
      this._endPixel = undefined;
      this._reversePixels = false;
    }

    init(options) {
      this.options = options;
      this.axis = options.axis;
    }

    setDimensions() {
      if (this.isHorizontal()) {
        this.width = this.maxWidth;
        this.left = 0;
        this.right = this.width;
      } else {
        this.height = this.maxHeight;
        this.top = 0;
        this.bottom = this.height;
      }
    }

    isHorizontal() {
      const {axis, position} = this.options;
      return position === 'top' || position === 'bottom' || axis === 'x';
    }

    getPixelForValue(value) {
      return this.getPixelForDecimal((value - this.min) / (this.max - this.min));
    }

    getPixelForDecimal(decimal) {
      if (this._reversePixels) decimal = 1 - decimal;
      const pixel = this._startPixel + decimal * (this._endPixel - this._startPixel);
      return Math.round(pixel);
    }

    configure() {
      let start, end;
      const reverse = this.options.reverse;
      
      if (this.isHorizontal()) {
        start = this.left;
        end = this.right;
      } else {
        start = this.top;
        end = this.bottom;
      }
      
      this._startPixel = start;
      this._endPixel = end;
      this._reversePixels = reverse;
    }

    update(maxWidth, maxHeight, margins) {
      this.maxWidth = maxWidth;
      this.maxHeight = maxHeight;
      this.setDimensions();
      this.configure();
    }
  }

  // Linear Scale (minimal)
  class LinearScale extends Scale {
    static id = 'linear';

    determineDataLimits() {
      const {min, max} = this.getMinMax(true);
      this.min = isFiniteNumber(min) ? min : 0;
      this.max = isFiniteNumber(max) ? max : 1;
      
      if (this.options.beginAtZero) {
        if (this.min > 0) this.min = 0;
        if (this.max < 0) this.max = 0;
      }
    }

    getMinMax() {
      // Simplified - just return reasonable defaults
      return {min: 0, max: 100};
    }

    buildTicks() {
      const opts = this.options;
      const tickCount = 6; // Fixed number of ticks
      const step = (this.max - this.min) / (tickCount - 1);
      const ticks = [];
      
      for (let i = 0; i < tickCount; i++) {
        ticks.push({
          value: this.min + (step * i)
        });
      }
      
      return ticks;
    }

    draw() {
      // Minimal implementation - no visual rendering needed for this dashboard
    }
  }

  // Dataset Controller base
  class DatasetController {
    constructor(chart, datasetIndex) {
      this.chart = chart;
      this.index = datasetIndex;
      this._cachedMeta = this.getMeta();
      this._type = this._cachedMeta.type;
      this.initialize();
    }

    initialize() {
      this.configure();
      this.linkScales();
      this.addElements();
    }

    configure() {
      // Basic configuration
    }

    linkScales() {
      const chart = this.chart;
      const meta = this._cachedMeta;
      const dataset = this.getDataset();
      
      meta.xScale = chart.scales.x;
      meta.yScale = chart.scales.y;
    }

    getDataset() {
      return this.chart.data.datasets[this.index];
    }

    getMeta() {
      return this.chart.getDatasetMeta(this.index);
    }

    addElements() {
      const meta = this._cachedMeta;
      meta.dataset = new LineElement();
      
      const data = this.getDataset().data;
      meta.data = data.map(() => new PointElement());
    }

    update(mode) {
      const meta = this._cachedMeta;
      const dataset = meta.dataset;
      const points = meta.data || [];
      const options = this.resolveDatasetElementOptions();
      
      // Update line element
      dataset.points = points;
      Object.assign(dataset, {
        options: options,
        _chart: this.chart,
        _datasetIndex: this.index
      });
      
      this.updateElements(points, 0, points.length, mode);
    }

    updateElements(points, start, count, mode) {
      const dataset = this.getDataset();
      const {xScale, yScale} = this._cachedMeta;
      
      for (let i = start; i < start + count; i++) {
        const point = points[i];
        const parsed = this.getParsed(i);
        
        if (parsed) {
          point.x = xScale.getPixelForValue(parsed.x);
          point.y = yScale.getPixelForValue(parsed.y);
          point.skip = false;
          point.options = this.resolveDataElementOptions(i);
        }
      }
    }

    resolveDatasetElementOptions() {
      const dataset = this.getDataset();
      return merge({}, dataset);
    }

    resolveDataElementOptions(index) {
      return this.resolveDatasetElementOptions();
    }

    getParsed(index) {
      const data = this.getDataset().data;
      if (index >= 0 && index < data.length) {
        return {x: index, y: data[index]};
      }
      return null;
    }

    draw() {
      const meta = this._cachedMeta;
      const dataset = meta.dataset;
      
      if (dataset) {
        dataset.draw(this.chart.ctx, this.chart.chartArea);
      }
    }
  }

  // Line Controller
  class LineController extends DatasetController {
    static id = 'line';
    static defaults = {
      datasetElementType: 'line',
      dataElementType: 'point',
      showLine: true,
      spanGaps: false
    };

    initialize() {
      this.enableOptionSharing = true;
      super.initialize();
    }
  }

  // Minimal Chart class
  class Chart {
    static version = '4.4.0-minimal';
    static instances = {};
    static registry = {
      controllers: new Map(),
      elements: new Map(),
      scales: new Map()
    };

    constructor(item, config) {
      const canvas = this._resolveCanvas(item);
      const ctx = canvas.getContext('2d');
      
      this.id = this._generateId();
      this.ctx = ctx;
      this.canvas = canvas;
      this.config = this._resolveConfig(config);
      this.data = this.config.data;
      this.options = this.config.options;
      this.width = canvas.width;
      this.height = canvas.height;
      this.scales = {};
      this.chartArea = {};
      this._metasets = [];
      
      this._initialize();
      Chart.instances[this.id] = this;
    }

    _generateId() {
      return Date.now() + Math.random().toString(36).substr(2, 9);
    }

    _resolveCanvas(item) {
      if (typeof item === 'string') {
        item = document.getElementById(item);
      }
      return item && item.canvas ? item.canvas : item;
    }

    _resolveConfig(config) {
      config = config || {};
      config.data = config.data || {datasets: [], labels: []};
      config.options = config.options || {};
      return config;
    }

    _initialize() {
      this.buildOrUpdateScales();
      this.buildOrUpdateControllers();
      this._updateLayout();
      this.bindEvents();
    }

    buildOrUpdateScales() {
      const options = this.options;
      const scalesConfig = options.scales || {};
      
      // Create default scales
      if (!scalesConfig.x) {
        scalesConfig.x = {type: 'category'};
      }
      if (!scalesConfig.y) {
        scalesConfig.y = {type: 'linear', beginAtZero: true};
      }
      
      Object.keys(scalesConfig).forEach(id => {
        const scaleConfig = scalesConfig[id];
        const scale = new LinearScale({id, type: scaleConfig.type, ctx: this.ctx, chart: this});
        scale.init(merge({axis: id}, scaleConfig));
        this.scales[id] = scale;
      });
    }

    buildOrUpdateControllers() {
      const newControllers = [];
      const datasets = this.data.datasets;
      
      for (let i = 0; i < datasets.length; i++) {
        const dataset = datasets[i];
        let meta = this.getDatasetMeta(i);
        const type = dataset.type || this.config.type;
        
        meta.type = type;
        meta.index = i;
        meta.visible = true;
        
        if (!meta.controller) {
          // Only support line charts
          meta.controller = new LineController(this, i);
          newControllers.push(meta.controller);
        }
      }
      
      return newControllers;
    }

    getDatasetMeta(datasetIndex) {
      const dataset = this.data.datasets[datasetIndex];
      const metasets = this._metasets;
      
      let meta = metasets.find(m => m && m._dataset === dataset);
      if (!meta) {
        meta = {
          type: null,
          data: [],
          dataset: null,
          controller: null,
          hidden: null,
          xAxisID: null,
          yAxisID: null,
          order: dataset && dataset.order || 0,
          index: datasetIndex,
          _dataset: dataset,
          _parsed: [],
          _sorted: false
        };
        metasets.push(meta);
      }
      
      return meta;
    }

    _updateLayout() {
      const canvas = this.canvas;
      const padding = 20; // Fixed padding
      
      this.chartArea = {
        left: padding,
        top: padding,
        right: canvas.width - padding,
        bottom: canvas.height - padding,
        width: canvas.width - (2 * padding),
        height: canvas.height - (2 * padding)
      };
      
      // Update scales dimensions
      Object.values(this.scales).forEach(scale => {
        if (scale.isHorizontal()) {
          scale.left = this.chartArea.left;
          scale.right = this.chartArea.right;
          scale.top = this.chartArea.bottom;
          scale.bottom = this.chartArea.bottom;
          scale.width = this.chartArea.width;
          scale.height = 0;
        } else {
          scale.left = this.chartArea.left;
          scale.right = this.chartArea.left;
          scale.top = this.chartArea.top;
          scale.bottom = this.chartArea.bottom;
          scale.width = 0;
          scale.height = this.chartArea.height;
        }
        scale.configure();
      });
    }

    update(mode) {
      // Update scales
      Object.values(this.scales).forEach(scale => {
        if (scale.determineDataLimits) {
          scale.determineDataLimits();
        }
        if (scale.buildTicks) {
          scale.ticks = scale.buildTicks();
        }
      });
      
      // Update controllers
      this._metasets.forEach(meta => {
        if (meta.controller) {
          meta.controller.update(mode);
        }
      });
      
      this.render();
    }

    render() {
      this.draw();
    }

    draw() {
      const ctx = this.ctx;
      const chartArea = this.chartArea;
      
      // Clear canvas
      ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
      
      if (this.width <= 0 || this.height <= 0) return;
      
      // Draw datasets
      this._metasets.forEach(meta => {
        if (meta.controller && meta.visible !== false) {
          meta.controller.draw();
        }
      });
    }

    bindEvents() {
      // Minimal event binding for resize
      const resizeHandler = () => {
        if (this.canvas.parentNode) {
          const container = this.canvas.parentNode;
          const containerWidth = container.clientWidth;
          const containerHeight = container.clientHeight;
          
          if (containerWidth !== this.width || containerHeight !== this.height) {
            this.resize(containerWidth, containerHeight);
          }
        }
      };
      
      window.addEventListener('resize', resizeHandler);
    }

    resize(width, height) {
      if (width) this.canvas.width = width;
      if (height) this.canvas.height = height;
      this.width = this.canvas.width;
      this.height = this.canvas.height;
      this._updateLayout();
      this.update();
    }

    destroy() {
      delete Chart.instances[this.id];
    }

    // Static method to register components
    static register(...items) {
      items.forEach(item => {
        if (item.id && item.defaults) {
          if (typeof item.prototype !== 'undefined' && item.prototype.constructor === item) {
            // It's a controller or scale
            if (item.id.includes('Scale') || item.prototype instanceof Scale) {
              Chart.registry.scales.set(item.id, item);
            } else if (item.prototype instanceof DatasetController) {
              Chart.registry.controllers.set(item.id, item);
            } else if (item.prototype instanceof Element) {
              Chart.registry.elements.set(item.id, item);
            }
          }
        }
      });
    }
  }

  // Register minimal components
  Chart.register(
    LineController,
    LinearScale,
    PointElement,
    LineElement
  );

  // Add essential properties to Chart
  Chart.Element = Element;
  Chart.DatasetController = DatasetController;
  Chart.Scale = Scale;

  return Chart;
});