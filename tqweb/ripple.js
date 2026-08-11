/**
 * Material-style click-origin ripple for [data-ripple] hosts.
 */
(() => {
  const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)');

  document.addEventListener('pointerdown', (event) => {
    if (reducedMotion.matches) {
      return;
    }

    const host = event.target.closest('[data-ripple]');
    if (!host) {
      return;
    }

    const rect = host.getBoundingClientRect();
    const size = Math.max(rect.width, rect.height);
    const x = event.clientX - rect.left - size / 2;
    const y = event.clientY - rect.top - size / 2;

    const wave = document.createElement('span');
    wave.className = 'ripple-wave';
    wave.style.width = `${size}px`;
    wave.style.height = `${size}px`;
    wave.style.left = `${x}px`;
    wave.style.top = `${y}px`;

    host.appendChild(wave);
    wave.addEventListener('animationend', () => wave.remove(), { once: true });
  });
})();
