window.addEventListener('DOMContentLoaded', () => {
    const container = document.getElementById('test-editormd');
    if (!container || typeof window.editormd !== 'function') {
        return;
    }
    window.goblogEditor = window.editormd('test-editormd', {
        width: '100%',
        height: 640,
        syncScrolling: 'single',
        path: container.dataset.editorPath
    });
});
