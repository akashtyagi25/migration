document.addEventListener('DOMContentLoaded', () => {
    const startBtn = document.getElementById('startBtn');
    const terminal = document.getElementById('terminal');
    const graphContainer = document.getElementById('graphContainer');
    
    // Stats
    const statFiles = document.getElementById('statFiles');
    const statPass = document.getElementById('statPass');
    const statFail = document.getElementById('statFail');
    const statAttempts = document.getElementById('statAttempts');

    let eventSource = null;

    function appendLog(message, type = 'info') {
        const div = document.createElement('div');
        
        // Color coding based on log type
        if (type === 'error' || message.includes('ERROR') || message.includes('FATAL')) {
            div.className = 'text-red-400';
        } else if (type === 'success' || message.includes('SUCCESS')) {
            div.className = 'text-green-400';
        } else if (message.includes('Attempt')) {
            div.className = 'text-purple-400';
        } else if (message.includes('MIGRATING:')) {
            div.className = 'text-yellow-400 font-bold mt-2';
        } else {
            div.className = 'text-gray-300';
        }
        
        // Ensure newlines are rendered
        div.innerHTML = message.replace(/\n/g, '<br>');
        terminal.appendChild(div);
        
        // Auto-scroll to bottom
        terminal.scrollTop = terminal.scrollHeight;
    }

    function createGraphNode(filename) {
        const div = document.createElement('div');
        div.id = `node-${filename.replace(/[^a-zA-Z0-9]/g, '-')}`;
        div.className = 'graph-node px-3 py-1.5 bg-gray-800 border border-gray-600 rounded-md text-xs font-mono text-gray-300 shadow-sm flex items-center space-x-2';
        div.innerHTML = `
            <svg class="w-3 h-3 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"></path></svg>
            <span>${filename}</span>
        `;
        return div;
    }

    function updateNodeStatus(filename, status) {
        const id = `node-${filename.replace(/[^a-zA-Z0-9]/g, '-')}`;
        const node = document.getElementById(id);
        if (node) {
            node.classList.remove('node-processing', 'node-success', 'node-error');
            if (status === 'processing') node.classList.add('node-processing');
            if (status === 'success') node.classList.add('node-success');
            if (status === 'error') node.classList.add('node-error');
        }
    }

    startBtn.addEventListener('click', () => {
        const targetLang = document.getElementById('targetLang').value;
        const instructions = document.getElementById('instructions').value;

        // Reset UI
        terminal.innerHTML = '';
        graphContainer.innerHTML = '';
        statFiles.innerText = '0';
        statPass.innerText = '0';
        statFail.innerText = '0';
        statAttempts.innerText = '0';
        
        appendLog(`> Starting migration to ${targetLang.toUpperCase()}...`, 'info');
        
        startBtn.disabled = true;
        startBtn.innerHTML = '<span class="animate-pulse">🔄 Running...</span>';

        // Initialize SSE Connection
        const url = `/api/migrate?target=${encodeURIComponent(targetLang)}&instructions=${encodeURIComponent(instructions)}`;
        
        if (eventSource) {
            eventSource.close();
        }

        eventSource = new EventSource(url);

        eventSource.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                
                if (data.type === 'log') {
                    appendLog(data.message);
                } else if (data.type === 'graph') {
                    // Update graph UI
                    graphContainer.innerHTML = '';
                    const files = data.files || [];
                    statFiles.innerText = files.length;
                    
                    files.forEach(file => {
                        graphContainer.appendChild(createGraphNode(file));
                    });
                } else if (data.type === 'file_start') {
                    updateNodeStatus(data.file, 'processing');
                } else if (data.type === 'file_success') {
                    updateNodeStatus(data.file, 'success');
                    statPass.innerText = parseInt(statPass.innerText) + 1;
                } else if (data.type === 'file_error') {
                    updateNodeStatus(data.file, 'error');
                    statFail.innerText = parseInt(statFail.innerText) + 1;
                } else if (data.type === 'attempt') {
                    statAttempts.innerText = parseInt(statAttempts.innerText) + 1;
                }

            } catch (e) {
                appendLog(`Parse error: ${e.message}`, 'error');
            }
        };

        eventSource.addEventListener('done', (event) => {
            appendLog('> Migration Batch Process Completed.', 'success');
            eventSource.close();
            startBtn.disabled = false;
            startBtn.innerHTML = '🚀 Start Migration';
        });

        eventSource.onerror = (err) => {
            console.error('SSE Error', err);
            appendLog('> Error connecting to migration server or connection closed.', 'error');
            eventSource.close();
            startBtn.disabled = false;
            startBtn.innerHTML = '🚀 Start Migration';
        };
    });
});
