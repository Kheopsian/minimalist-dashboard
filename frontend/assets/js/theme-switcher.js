document.addEventListener("DOMContentLoaded", () => {
  const themeSwitcher = document.getElementById("theme-switcher");
  const body = document.body;
  const themeIcon = themeSwitcher.querySelector("i");

  const savedTheme = localStorage.getItem("theme");
  if (savedTheme === "dark") {
    body.classList.add("dark-mode");
    themeIcon.classList.remove("fa-moon");
    themeIcon.classList.add("fa-sun");
  }

  themeSwitcher.addEventListener("click", () => {
    body.classList.toggle("dark-mode");

    if (body.classList.contains("dark-mode")) {
      themeIcon.classList.remove("fa-moon");
      themeIcon.classList.add("fa-sun");
      localStorage.setItem("theme", "dark");
    } else {
      themeIcon.classList.remove("fa-sun");
      themeIcon.classList.add("fa-moon");
      localStorage.setItem("theme", "light");
    }
  });
});
