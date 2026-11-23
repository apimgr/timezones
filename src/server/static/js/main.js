// Theme cycle: dark -> light -> auto -> dark.
const THEME_CYCLE = ['dark', 'light', 'auto'];

function currentTheme() {
  const match = document.documentElement.className.match(/theme-(dark|light|auto)/);
  return match ? match[1] : 'dark';
}

function setTheme(theme) {
  document.documentElement.className = `theme-${theme}`;
  document.cookie = `theme=${theme}; path=/; max-age=31536000; SameSite=Lax`;
  const button = document.querySelector('.theme-toggle');
  const icon = button?.querySelector('.theme-icon');
  const icons = { dark: '🌙', light: '☀️', auto: '◐' };
  if (icon) icon.textContent = icons[theme];
  if (button) {
    button.setAttribute('aria-label', `Switch theme (currently ${theme})`);
    button.title = `Switch theme (currently ${theme})`;
  }
}

document.querySelectorAll('.theme-toggle-form').forEach((form) => {
  form.addEventListener('submit', (event) => {
    event.preventDefault();
    const next = THEME_CYCLE[(THEME_CYCLE.indexOf(currentTheme()) + 1) % THEME_CYCLE.length];
    setTheme(next);
  });
});
