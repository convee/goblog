(() => {
    "use strict";

    document.querySelectorAll('form[method="GET"], form[method="get"]').forEach((form) => {
        form.addEventListener("formdata", (event) => {
            for (const [key, value] of event.formData.entries()) {
                if (typeof value === "string" && value.trim() === "") {
                    event.formData.delete(key);
                }
            }
            if (event.formData.get("per_page") === "20") {
                event.formData.delete("per_page");
            }
        });
    });

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
