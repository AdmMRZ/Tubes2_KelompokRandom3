document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('search-form');
  const btnText = document.querySelector('.btn-text');
  const btnLoader = document.querySelector('.btn-loader');
  const submitBtn = document.getElementById('submit-btn');
  
  const resultsArea = document.getElementById('results-area');
  
  const metricTime = document.getElementById('metric-time');
  const metricEvaluated = document.getElementById('metric-evaluated');
  const metricFound = document.getElementById('metric-found');
  const metricDepth = document.getElementById('metric-depth');
  const algoUsed = document.getElementById('algo-used');
  
  const traversalContainer = document.getElementById('traversal-container');
  const matchList = document.getElementById('match-list');
  const treeContainer = document.getElementById('tree-container');

  const modeRadios = document.querySelectorAll('input[name="input_mode"]');
  const urlGroup = document.getElementById('url-group');
  const htmlGroup = document.getElementById('html-group');

  modeRadios.forEach(radio => {
    radio.addEventListener('change', (e) => {
      if (e.target.value === 'url') {
        urlGroup.classList.remove('hidden');
        htmlGroup.classList.add('hidden');
      } else {
        urlGroup.classList.add('hidden');
        htmlGroup.classList.remove('hidden');
      }
    });
  });

  form.addEventListener('submit', async (e) => {
    e.preventDefault();

    const mode = document.querySelector('input[name="input_mode"]:checked').value;
    const url = document.getElementById('url').value;
    const html_content = document.getElementById('html-text').value;
    const selector = document.getElementById('selector').value;
    const algo = document.getElementById('algo').value;
    const limit = parseInt(document.getElementById('limit').value) || 0;

    btnText.textContent = 'Processing...';
    btnLoader.classList.remove('hidden');
    submitBtn.disabled = true;
    
    resultsArea.classList.add('hidden');
    traversalContainer.innerHTML = '';
    matchList.innerHTML = '';
    treeContainer.innerHTML = '';

    try {
      const payload = {
        selector,
        algo,
        limit
      };

      if (mode === 'url') {
        if (!url) throw new Error("Target URL is required in URL Mode.");
        payload.url = url;
      } else {
        if (!html_content.trim()) throw new Error("Raw HTML is required in Text Mode.");
        payload.html_content = html_content;
      }

      const response = await fetch('http://localhost:8080/api/search', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
      });

      if (!response.ok) {
        const errText = await response.text();
        throw new Error(errText || `HTTP Error ${response.status}`);
      }

      const data = await response.json();
      
      metricTime.textContent = data.execution_time_ms;
      metricEvaluated.textContent = data.node_count;
      metricFound.textContent = data.results ? data.results.length : 0;
      metricDepth.textContent = data.max_depth || 0;
      algoUsed.textContent = algo;

      resultsArea.classList.remove('hidden');

      animateTraversalLog(data.traversal_log, selector);

      renderMatches(data.results);

      if (data.root_tree) {
        renderTree(data.root_tree, treeContainer);
      }

    } catch (error) {
      alert(`An error occurred:\n${error.message}`);
    } finally {
      btnText.textContent = 'Execute Traversal';
      btnLoader.classList.add('hidden');
      submitBtn.disabled = false;
    }
  });

  function getAttrsHtml(attrs) {
    if (!attrs) return '';
    const keys = Object.keys(attrs);
    let html = '';
    keys.forEach(k => {
      html += ` <span class="tree-attr-key">${k}</span>="<span class="tree-attr-val">${attrs[k]}</span>"`;
    });
    return html;
  }

  function renderTree(node, container) {
    if (!node) return;
    
    if (!node.Children || node.Children.length === 0) {
      const div = document.createElement('div');
      div.className = 'tree-leaf';
      div.innerHTML = `&lt;${node.Tag}<span class="tree-attrs">${getAttrsHtml(node.Attributes)}</span>&gt;`;
      container.appendChild(div);
      return;
    }

    const details = document.createElement('details');
    if (container === treeContainer) details.open = true;

    const summary = document.createElement('summary');
    summary.innerHTML = `&lt;${node.Tag}<span class="tree-attrs">${getAttrsHtml(node.Attributes)}</span>&gt;`;
    
    details.appendChild(summary);
    container.appendChild(details);

    node.Children.forEach(child => renderTree(child, details));
  }

  function animateTraversalLog(logs, selectorTarget) {
    if (!logs) return;
    
    let i = 0;
    const chunkSize = 50; 
    
    function drawChunk() {
      const end = Math.min(i + chunkSize, logs.length);
      const fragment = document.createDocumentFragment();
      
      for (; i < end; i++) {
        const item = logs[i];
        const span = document.createElement('span');
        span.className = 'chip';
        span.textContent = item;
        
        span.style.animationDelay = `${(i % chunkSize) * 10}ms`;
        fragment.appendChild(span);
      }
      
      traversalContainer.appendChild(fragment);
      traversalContainer.scrollTop = traversalContainer.scrollHeight;

      if (i < logs.length) {
        requestAnimationFrame(drawChunk);
      }
    }
    
    requestAnimationFrame(drawChunk);
  }

  function renderMatches(results) {
    if (!results || results.length === 0) {
      const li = document.createElement('li');
      li.className = 'match-item';
      li.style.textAlign = 'center';
      li.textContent = 'No matching nodes found.';
      matchList.appendChild(li);
      return;
    }

    const dispLimit = 500;
    const total = results.length;
    
    results.slice(0, dispLimit).forEach(node => {
      const li = document.createElement('li');
      li.className = 'match-item';
      
      li.innerHTML = `
        <div class="tag">&lt;${node.Tag}&gt;</div>
        <div class="attrs">${getAttrsHtml(node.Attributes)}</div>
        <div class="children">Contains ${node.Children ? node.Children.length : 0} direct children</div>
      `;
      matchList.appendChild(li);
    });

    if (total > dispLimit) {
      const li = document.createElement('li');
      li.className = 'match-item';
      li.style.textAlign = 'center';
      li.textContent = `... and ${total - dispLimit} more matches.`;
      matchList.appendChild(li);
    }
  }
});
