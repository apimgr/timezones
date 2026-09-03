// Timezones API - Main JavaScript

// Theme toggle
function toggleTheme() {
  const currentTheme = document.documentElement.getAttribute('data-theme');
  const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
  document.documentElement.setAttribute('data-theme', newTheme);
  localStorage.setItem('theme', newTheme);

  // Update icon
  const icon = document.querySelector('.theme-icon');
  icon.textContent = newTheme === 'dark' ? '🌙' : '☀️';
}

// Initialize theme from localStorage
document.addEventListener('DOMContentLoaded', () => {
  const savedTheme = localStorage.getItem('theme') || 'dark';
  document.documentElement.setAttribute('data-theme', savedTheme);

  const icon = document.querySelector('.theme-icon');
  if (icon) {
    icon.textContent = savedTheme === 'dark' ? '🌙' : '☀️';
  }
});
