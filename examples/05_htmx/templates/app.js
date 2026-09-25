// Stands in for the HTMX server round-trip: keeps the hidden-series list,
// asks the Go WASM module for the re-rendered fragment and swaps it in.
const container = document.getElementById('chart-container');
const hidden = new Set();

function render() {
    container.innerHTML = goRender([...hidden].join(','));
}

container.addEventListener('click', (ev) => {
    const btn = ev.target.closest('button[data-series]');
    if (!btn) return;
    const name = btn.dataset.series;
    if (hidden.has(name)) hidden.delete(name); else hidden.add(name);
    render();
});

window.wasmReady = render;

async function loadWASM() {
    try {
        const go = new Go();
        const result = await WebAssembly.instantiateStreaming(
            fetch('main.wasm'), go.importObject
        );
        go.run(result.instance);
    } catch (err) {
        container.innerHTML = '<div class="notification is-danger">Failed to load WASM: ' + err.message + '</div>';
    }
}

loadWASM();
