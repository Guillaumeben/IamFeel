/**
 * Toast Notification System
 * Provides non-intrusive feedback for user actions
 */

class ToastManager {
    constructor() {
        this.container = null;
        this.toasts = new Map();
        this.init();
    }

    init() {
        // Create toast container if it doesn't exist
        if (!document.getElementById('toast-container')) {
            this.container = document.createElement('div');
            this.container.id = 'toast-container';
            this.container.className = 'toast-container';
            document.body.appendChild(this.container);
        } else {
            this.container = document.getElementById('toast-container');
        }
    }

    /**
     * Show a toast notification
     * @param {string} type - Type of toast: 'success', 'error', 'info', 'warning'
     * @param {string} title - Toast title
     * @param {string} message - Toast message
     * @param {number} duration - Duration in ms (0 for permanent)
     */
    show(type, title, message, duration = 4000) {
        const id = `toast-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

        const icons = {
            success: '✓',
            error: '⚠️',
            info: 'ℹ️',
            warning: '⚡'
        };

        const toast = document.createElement('div');
        toast.id = id;
        toast.className = `toast toast-${type}`;
        toast.innerHTML = `
            <div class="toast-icon">${icons[type] || 'ℹ️'}</div>
            <div class="toast-content">
                <div class="toast-title">${title}</div>
                ${message ? `<div class="toast-message">${message}</div>` : ''}
            </div>
            <button class="toast-close" aria-label="Close notification">&times;</button>
        `;

        // Add to container
        this.container.appendChild(toast);
        this.toasts.set(id, toast);

        // Close button handler
        const closeBtn = toast.querySelector('.toast-close');
        closeBtn.addEventListener('click', () => this.hide(id));

        // Auto-hide after duration
        if (duration > 0) {
            setTimeout(() => this.hide(id), duration);
        }

        return id;
    }

    /**
     * Hide a toast notification
     * @param {string} id - Toast ID
     */
    hide(id) {
        const toast = this.toasts.get(id);
        if (!toast) return;

        toast.classList.add('toast-hiding');

        setTimeout(() => {
            if (toast.parentNode) {
                toast.parentNode.removeChild(toast);
            }
            this.toasts.delete(id);
        }, 300); // Match animation duration
    }

    /**
     * Clear all toasts
     */
    clearAll() {
        this.toasts.forEach((_, id) => this.hide(id));
    }

    // Convenience methods
    success(title, message, duration) {
        return this.show('success', title, message, duration);
    }

    error(title, message, duration) {
        return this.show('error', title, message, duration);
    }

    info(title, message, duration) {
        return this.show('info', title, message, duration);
    }

    warning(title, message, duration) {
        return this.show('warning', title, message, duration);
    }
}

// Create global toast instance
window.toast = new ToastManager();

// Check for URL parameters to show toasts (for redirects)
document.addEventListener('DOMContentLoaded', () => {
    const urlParams = new URLSearchParams(window.location.search);

    if (urlParams.has('success')) {
        toast.success('Success', urlParams.get('success'));
        // Clean URL
        const cleanUrl = window.location.pathname;
        window.history.replaceState({}, document.title, cleanUrl);
    }

    if (urlParams.has('error')) {
        toast.error('Error', urlParams.get('error'));
        // Clean URL
        const cleanUrl = window.location.pathname;
        window.history.replaceState({}, document.title, cleanUrl);
    }

    if (urlParams.has('info')) {
        toast.info('Info', urlParams.get('info'));
        // Clean URL
        const cleanUrl = window.location.pathname;
        window.history.replaceState({}, document.title, cleanUrl);
    }
});
