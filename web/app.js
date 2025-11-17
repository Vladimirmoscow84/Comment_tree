
// GET  /comments?limit=..&offset=..&search=..&parent=..&sort=asc|desc
// POST /comments { parent_id, content }
// DELETE /comments/:id

const treeEl = document.getElementById('tree');
const searchInput = document.getElementById('searchInput');
const searchBtn = document.getElementById('searchBtn');
const postRootBtn = document.getElementById('postRoot');
const rootContent = document.getElementById('rootContent');
const sortSelect = document.getElementById('sortSelect');

const prevPage = document.getElementById('prevPage');
const nextPage = document.getElementById('nextPage');
const prevPage2 = document.getElementById('prevPage2');
const nextPage2 = document.getElementById('nextPage2');
const pageInfo = document.getElementById('pageInfo');
const pageInfo2 = document.getElementById('pageInfo2');

let page = 1;
const perPage = 10;
let totalLastFetched = 0;


const subtreeCache = new Map();

function apiGet(url){
  return fetch(url).then(res => {
    if(!res.ok) throw new Error('HTTP ' + res.status);
    return res.json();
  });
}

function apiPost(url, body){
  return fetch(url, {
    method: 'POST',
    headers: {'Content-Type':'application/json'},
    body: JSON.stringify(body)
  }).then(res => {
    if(!res.ok) return res.json().then(j => { throw new Error(j.error || JSON.stringify(j)) });
    return res.json();
  });
}

function apiDelete(url){
  return fetch(url, { method: 'DELETE' }).then(res => {
    if(!res.ok) throw new Error('HTTP ' + res.status);
    return res.json().catch(()=>({}));
  });
}

function renderTree(flatList){
 
  const map = new Map();
  flatList.forEach(c => { c.children = []; map.set(c.id, c); });
  const roots = [];
  flatList.forEach(c => {
    if(c.parent_id){
      const p = map.get(c.parent_id);
      if(p) p.children.push(c);
      else roots.push(c); // orphan
    } else {
      roots.push(c);
    }
  });

  treeEl.innerHTML = '';
  roots.forEach(node => renderNode(node, 0, treeEl));
}

function renderNode(node, depth, container){
  const wrap = document.createElement('div');
  wrap.className = 'comment';
  wrap.style.marginLeft = (depth * 18) + 'px';

  const meta = document.createElement('div');
  meta.className = 'meta';
  meta.innerHTML = `#${node.id} • ${new Date(node.created_at).toLocaleString()}`;
  wrap.appendChild(meta);

  const content = document.createElement('div');
  content.className = 'content';
  content.textContent = node.content;
  wrap.appendChild(content);

  const actions = document.createElement('div');
  actions.className = 'actions';

 
  const expandBtn = document.createElement('button');
  expandBtn.className = 'btn btn-expand';
  expandBtn.textContent = 'Развернуть';
  expandBtn.onclick = () => onToggleChildren(node, wrap, depth);
  actions.appendChild(expandBtn);

  const replyBtn = document.createElement('button');
  replyBtn.className = 'btn btn-reply';
  replyBtn.textContent = 'Ответить';
  replyBtn.onclick = () => onReply(node, wrap, depth);
  actions.appendChild(replyBtn);

  const delBtn = document.createElement('button');
  delBtn.className = 'btn btn-delete';
  delBtn.textContent = 'Удалить';
  delBtn.onclick = async () => {
    if(!confirm('Удалить комментарий и все ответы?')) return;
    try {
      await apiDelete(`/comments/${node.id}`);
      refresh(); // reload current view
    } catch (err) {
      alert('Ошибка удаления: ' + err.message);
    }
  };
  actions.appendChild(delBtn);

  wrap.appendChild(actions);

  container.appendChild(wrap);

 
  if(subtreeCache.has(node.id)){
    const children = subtreeCache.get(node.id);
    if(children && children.length){
      const inner = document.createElement('div');
      inner.className = 'indent';
      children.forEach(ch => renderNode(ch, depth+1, inner));
      container.appendChild(inner);
    }
  }
}

async function onToggleChildren(node, el, depth){

  if(subtreeCache.has(node.id)){
    subtreeCache.delete(node.id);
    refresh(); 
    return;
  }

  
  try {
    const list = await apiGet(`/comments?parent=${encodeURIComponent(node.id)}`);
    const children = list.filter(c => c.id !== node.id);
    subtreeCache.set(node.id, children);
    refresh(); 
  } catch (err) {
    alert('Ошибка загрузки ветки: ' + err.message);
  }
}

function onReply(node, el, depth){
  const promptText = prompt('Текст ответа:');
  if(!promptText) return;
  apiPost('/comments', { parent_id: node.id, content: promptText })
    .then(()=> refresh())
    .catch(err => alert('Ошибка создания ответа: ' + err.message));
}

async function loadRoots(opts = {}){
  const limit = opts.limit ?? perPage;
  const offset = ((page - 1) * perPage);
  const sort = sortSelect.value || 'desc';
  const url = `/comments?limit=${limit}&offset=${offset}&sort=${encodeURIComponent(sort)}`;
  const list = await apiGet(url);
  totalLastFetched = list.length;
  return list;
}

async function onSearch(){
  const q = searchInput.value.trim();
  if(!q){
    page = 1;
    subtreeCache.clear();
    await refresh();
    return;
  }
  const list = await apiGet(`/comments?search=${encodeURIComponent(q)}&limit=100&offset=0`);
  treeEl.innerHTML = '';
  list.forEach(c => {
    const div = document.createElement('div');
    div.className = 'comment';
    div.innerHTML = `<div class="meta">#${c.id} • ${new Date(c.created_at).toLocaleString()}</div>
                     <div class="content">${escapeHtml(c.content)}</div>
                     <div class="actions"><button class="btn btn-reply" data-id="${c.id}">Ответить</button>
                     <button class="btn btn-delete" data-id="${c.id}">Удалить</button></div>`;
    treeEl.appendChild(div);
  });

  treeEl.querySelectorAll('.btn-reply').forEach(btn => {
    btn.addEventListener('click', (e) => {
      const id = Number(e.currentTarget.dataset.id);
      const text = prompt('Текст ответа:');
      if(!text) return;
      apiPost('/comments', { parent_id: id, content: text }).then(()=> onSearch()).catch(err=>alert(err.message));
    });
  });
  treeEl.querySelectorAll('.btn-delete').forEach(btn => {
    btn.addEventListener('click', (e) => {
      const id = Number(e.currentTarget.dataset.id);
      if(!confirm('Удалить?')) return;
      apiDelete(`/comments/${id}`).then(()=> onSearch()).catch(err=>alert(err.message));
    });
  });
}

function escapeHtml(s){
  return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
}

async function refresh(){
  const list = await loadRoots();
  renderTree(list);
  updatePagination();
}

function updatePagination(){
  pageInfo.textContent = `Страница ${page}`;
  pageInfo2.textContent = `Страница ${page}`;
  prevPage.disabled = (page <= 1);
  prevPage2.disabled = (page <= 1);
  nextPage.disabled = (totalLastFetched < perPage);
  nextPage2.disabled = (totalLastFetched < perPage);
}

prevPage.addEventListener('click', ()=>{ if(page>1){ page--; refresh(); }});
prevPage2.addEventListener('click', ()=>{ if(page>1){ page--; refresh(); }});
nextPage.addEventListener('click', ()=>{ page++; refresh(); });
nextPage2.addEventListener('click', ()=>{ page++; refresh(); });

searchBtn.addEventListener('click', ()=> onSearch());
sortSelect.addEventListener('change', ()=> refresh());

postRootBtn.addEventListener('click', async ()=>{
  const content = rootContent.value.trim();
  if(!content) return alert('Введите текст комментария');
  try {
    await apiPost('/comments', { parent_id: null, content });
    rootContent.value = '';
    subtreeCache.clear();
    page = 1;
    await refresh();
  } catch (err) {
    alert('Ошибка публикации: ' + err.message);
  }
});


refresh();
