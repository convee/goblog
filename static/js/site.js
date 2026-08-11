(() => {
    "use strict";

    const toggle = document.querySelector(".navbar-toggle[data-target]");
    if (!toggle) return;

    const target = document.querySelector(toggle.dataset.target);
    if (!target) return;

    toggle.addEventListener("click", () => {
        const expanded = toggle.getAttribute("aria-expanded") === "true";
        toggle.setAttribute("aria-expanded", String(!expanded));
        toggle.classList.toggle("collapsed", expanded);
        target.classList.toggle("in", !expanded);
    });
})();
