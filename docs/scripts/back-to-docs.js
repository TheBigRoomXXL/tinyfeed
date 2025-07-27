
function BackToDocumentation() {
    const referer = document.referrer;
    console.log(window.location.origin);
    console.log(document.referrer);
    if (referer.startsWith(window.location.origin)) {
        window.location.href = referer;
        return;
    }

    window.location.href = window.location.origin;
};

const button = document.createElement('button');
button.innerText = "Back to the documentation ↩";
button.onclick = BackToDocumentation;
button.style.margin = '1em';
button.style.float = "right";

const header = document.querySelector("header");
header?.appendChild(button);

