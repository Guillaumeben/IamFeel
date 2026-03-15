/**
 * Workout Details Renderer
 * Converts markdown-style workout text into structured HTML
 */

function renderWorkoutDetails(element) {
    const text = element.textContent;

    // Extract Focus (first line after "Focus:")
    const focusMatch = text.match(/Focus:\s*([^\n]+)/);

    if (!focusMatch) {
        // No focus found, keep original
        return;
    }

    let html = '';

    // Render Focus section
    const focus = focusMatch[1].trim();
    html += `<div class="workout-focus">
        <div class="workout-focus-label">Focus</div>
        <div class="workout-focus-text">${escapeHtml(focus)}</div>
    </div>`;

    // Get everything after the Focus line
    const focusIndex = text.indexOf('Focus:');
    const focusLineEnd = text.indexOf('\n', focusIndex);

    if (focusLineEnd !== -1) {
        const remainingText = text.substring(focusLineEnd + 1).trim();

        if (remainingText) {
            const detailsId = 'details-' + Math.random().toString(36).substr(2, 9);
            html += `<div class="workout-details-toggle">
                <button class="btn-toggle-details" onclick="toggleWorkoutDetails('${detailsId}')">
                    <span class="toggle-icon">▶</span>
                    <span class="toggle-text">Show Workout Details</span>
                </button>
            </div>
            <div id="${detailsId}" class="workout-details-container" style="display: none;">
                <div class="workout-details-simple">${formatPlainText(remainingText)}</div>
            </div>`;
        }
    }

    element.innerHTML = html;
    element.classList.add('workout-rendered');
}

// Simple text formatter - just preserve line breaks and basic formatting
function formatPlainText(text) {
    // Split into lines
    const lines = text.split('\n');
    let html = '';

    for (let line of lines) {
        line = line.trim();

        if (!line) {
            // Empty line = spacing
            html += '<div class="workout-spacer"></div>';
            continue;
        }

        // Bold headers (text between **)
        line = line.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');

        // Escape any remaining HTML
        const div = document.createElement('div');
        div.innerHTML = line;
        line = div.innerHTML;

        // Add the line
        html += `<div class="workout-line">${line}</div>`;
    }

    return html;
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Toggle workout details visibility
function toggleWorkoutDetails(detailsId) {
    const detailsContainer = document.getElementById(detailsId);
    const button = event.currentTarget;
    const icon = button.querySelector('.toggle-icon');
    const text = button.querySelector('.toggle-text');

    if (detailsContainer.style.display === 'none') {
        detailsContainer.style.display = 'block';
        icon.textContent = '▼';
        text.textContent = 'Hide Workout Details';
    } else {
        detailsContainer.style.display = 'none';
        icon.textContent = '▶';
        text.textContent = 'Show Workout Details';
    }
}

// Auto-render on page load
document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('.planned-notes, .session-notes').forEach(element => {
        if (!element.classList.contains('workout-rendered')) {
            renderWorkoutDetails(element);
        }
    });
});
