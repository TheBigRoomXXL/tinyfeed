
function onScroll() {
    if (document.body.scrollTop > 20 || document.documentElement.scrollTop > 20) {
        button.style.display = "block";
    } else {
        button.style.display = "none";
    }
}

function GoToTop() {
    document.body.scrollTop = 0; // For Safari
    document.documentElement.scrollTop = 0; // For Chrome, Firefox, IE and Opera
}


const button = document.createElement('button');
button.title = "Go to the top of the page";
button.onclick = GoToTop;
button.style.display = 'none';
button.style.position = 'fixed';
button.style.bottom = '30px';
button.style.right = '30px';
button.style.zIndex = '99';
button.style.cursor = 'pointer';
button.style.width = '3em';
button.style.height = '3em';
button.style.backgroundColor = 'current';
button.style.mask = 'size: 100 % 100 %';
button.style.maskRepeat = 'no';
button.style.maskImage = `url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23000' d='M22 12a10 10 0 0 1-10 10A10 10 0 0 1 2 12A10 10 0 0 1 12 2a10 10 0 0 1 10 10M7.4 15.4l4.6-4.6l4.6 4.6L18 14l-6-6l-6 6z'/%3E%3C/svg%3E")`;
button.style.webkitMask = 'size: 100 % 100 %';
button.style.webkitMaskRepeat = 'no';
button.style.webkitMaskImage = `url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23000' d='M22 12a10 10 0 0 1-10 10A10 10 0 0 1 2 12A10 10 0 0 1 12 2a10 10 0 0 1 10 10M7.4 15.4l4.6-4.6l4.6 4.6L18 14l-6-6l-6 6z'/%3E%3C/svg%3E")`;

document.body.appendChild(button);
window.onscroll = onScroll;
