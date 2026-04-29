(function(){let e=document.createElement(`link`).relList;if(e&&e.supports&&e.supports(`modulepreload`))return;for(let e of document.querySelectorAll(`link[rel="modulepreload"]`))n(e);new MutationObserver(e=>{for(let t of e)if(t.type===`childList`)for(let e of t.addedNodes)e.tagName===`LINK`&&e.rel===`modulepreload`&&n(e)}).observe(document,{childList:!0,subtree:!0});function t(e){let t={};return e.integrity&&(t.integrity=e.integrity),e.referrerPolicy&&(t.referrerPolicy=e.referrerPolicy),e.crossOrigin===`use-credentials`?t.credentials=`include`:e.crossOrigin===`anonymous`?t.credentials=`omit`:t.credentials=`same-origin`,t}function n(e){if(e.ep)return;e.ep=!0;let n=t(e);fetch(e.href,n)}})();var e=class{constructor(e){this.routes=e,window.addEventListener(`hashchange`,()=>this.resolve())}resolve(){let e=window.location.hash.slice(1)||`/login`;for(let[t,n]of Object.entries(this.routes)){let r=t.match(/^(.+)\/:(\w+)$/);if(r){let[,t,i]=r;if(e.startsWith(t+`/`)){let r=e.slice(t.length+1);if(r){n({[i]:r});return}}}}let t=this.routes[e];t?t():window.location.hash=`#/login`}start(){this.resolve()}},t=window.INSIGHTA_API_URL||`https://stage1.doxantro.com`;function n(){let e=document.cookie.match(/csrf_token=([^;]+)/);return e?e[1]:``}async function r(e,r={}){let a=`${t}${e}`,o={"X-API-Version":`1`,...r.headers};[`POST`,`PUT`,`DELETE`,`PATCH`].includes((r.method||`GET`).toUpperCase())&&(o[`X-CSRF-Token`]=n());let s=await fetch(a,{...r,headers:o,credentials:`include`});if(s.status===401){if(await i())return fetch(a,{...r,headers:o,credentials:`include`});throw window.location.hash=`#/login`,Error(`Session expired. Please log in again.`)}return s}async function i(){try{return(await fetch(`${t}/auth/refresh`,{method:`POST`,credentials:`include`,headers:{"Content-Type":`application/json`},body:JSON.stringify({})})).ok}catch{return!1}}function a(){return t}var o=`<svg viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>`;function s(e){let t=a();e.innerHTML=`
    <div class="login-page">
      <div class="login-card">
        <h1>Insighta <span style="color: var(--accent)">Labs+</span></h1>
        <p>Profile Intelligence Platform<br/>Sign in to access your dashboard</p>
        <button class="btn-github" id="github-login-btn">
          ${o}
          Continue with GitHub
        </button>
      </div>
    </div>
  `,document.getElementById(`github-login-btn`).addEventListener(`click`,()=>{window.location.href=`${t}/auth/github?client=web`})}async function c(e){w(e,`dashboard`,`
    <div class="page-header">
      <h1>Dashboard</h1>
      <p>Overview of the Profile Intelligence System</p>
    </div>
    <div class="stats-grid" id="stats-grid">
      <div class="loading"><div class="spinner"></div> Loading stats...</div>
    </div>
  `);try{let e=await(await r(`/api/profiles?page=1&limit=1`)).json();if(e.status!==`success`){document.getElementById(`stats-grid`).innerHTML=`<p>Failed to load stats</p>`;return}let t=(await(await r(`/api/profiles?page=1&limit=50`)).json()).data||[],n=t.filter(e=>e.gender===`male`).length,i=t.filter(e=>e.gender===`female`).length,a=[...new Set(t.map(e=>e.country_id))],o=t.length>0?Math.round(t.reduce((e,t)=>e+(t.age||0),0)/t.length):0;document.getElementById(`stats-grid`).innerHTML=`
      <div class="stat-card">
        <div class="stat-label">Total Profiles</div>
        <div class="stat-value">${e.total.toLocaleString()}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Male (sample)</div>
        <div class="stat-value" style="color: #60a5fa">${n}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Female (sample)</div>
        <div class="stat-value" style="color: #f472b6">${i}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Countries (sample)</div>
        <div class="stat-value">${a.length}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Avg Age (sample)</div>
        <div class="stat-value">${o}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Total Pages</div>
        <div class="stat-value">${e.total_pages}</div>
      </div>
    `}catch(e){document.getElementById(`stats-grid`).innerHTML=`<p>Error: ${e.message}</p>`}}var l=1,u=10,d={};async function f(e){w(e,`profiles`,`
    <div class="page-header">
      <h1>Profiles</h1>
      <p>Browse and filter the intelligence database</p>
    </div>

    <div class="table-container">
      <div class="table-toolbar">
        <select id="filter-gender" class="filter-select">
          <option value="">All Genders</option>
          <option value="male">Male</option>
          <option value="female">Female</option>
        </select>
        <input type="text" id="filter-country" class="filter-input" placeholder="Country Code (e.g. NG)" />
        <select id="filter-age" class="filter-select">
          <option value="">All Ages</option>
          <option value="child">Child</option>
          <option value="teenager">Teenager</option>
          <option value="adult">Adult</option>
          <option value="senior">Senior</option>
        </select>
        <button id="btn-apply-filters" class="btn-page">Apply</button>
        <button id="btn-export" class="btn-page" style="margin-left: auto;">Export CSV</button>
      </div>

      <div style="overflow-x: auto;">
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Gender</th>
              <th>Age</th>
              <th>Group</th>
              <th>Country</th>
              <th>Created</th>
            </tr>
          </thead>
          <tbody id="profiles-tbody">
            <tr><td colspan="6" class="loading"><div class="spinner"></div> Loading profiles...</td></tr>
          </tbody>
        </table>
      </div>

      <div class="pagination">
        <div class="pagination-info" id="pagination-info">Showing 0 results</div>
        <div class="pagination-buttons">
          <button id="btn-prev" class="btn-page" disabled>Previous</button>
          <button id="btn-next" class="btn-page" disabled>Next</button>
        </div>
      </div>
    </div>
  `),document.getElementById(`btn-apply-filters`).addEventListener(`click`,()=>{d.gender=document.getElementById(`filter-gender`).value,d.country_id=document.getElementById(`filter-country`).value,d.age_group=document.getElementById(`filter-age`).value,l=1,p()}),document.getElementById(`btn-export`).addEventListener(`click`,()=>{let e=new URLSearchParams({format:`csv`,...d});window.location.href=`${a()}/api/profiles/export?${e}`}),document.getElementById(`btn-prev`).addEventListener(`click`,()=>{l>1&&(l--,p())}),document.getElementById(`btn-next`).addEventListener(`click`,()=>{l++,p()}),document.getElementById(`filter-gender`).value=d.gender||``,document.getElementById(`filter-country`).value=d.country_id||``,document.getElementById(`filter-age`).value=d.age_group||``,p()}async function p(){let e=document.getElementById(`profiles-tbody`);e.innerHTML=`<tr><td colspan="6" class="loading"><div class="spinner"></div> Loading profiles...</td></tr>`;try{let t=await(await r(`/api/profiles?${new URLSearchParams({page:l,limit:u,...d})}`)).json();if(t.status!==`success`)throw Error(t.message||`Failed to load`);if(!t.data||t.data.length===0){e.innerHTML=`<tr><td colspan="6" style="padding: 32px; text-align: center; color: var(--text-muted);">No profiles found</td></tr>`,m(t);return}e.innerHTML=t.data.map(e=>`
      <tr onclick="window.location.hash='#/profiles/${e.id}'">
        <td style="font-weight: 500;">${h(e.name)}</td>
        <td><span class="badge badge-${e.gender}">${e.gender}</span></td>
        <td>${e.age}</td>
        <td style="text-transform: capitalize;">${e.age_group}</td>
        <td>${h(e.country_name)}</td>
        <td style="color: var(--text-muted); font-size: 0.8rem;">${new Date(e.created_at).toLocaleDateString()}</td>
      </tr>
    `).join(``),m(t)}catch(t){e.innerHTML=`<tr><td colspan="6" style="padding: 24px; color: var(--danger);">${h(t.message)}</td></tr>`}}function m(e){document.getElementById(`pagination-info`).textContent=`Page ${e.page} of ${e.total_pages} (Total: ${e.total})`,document.getElementById(`btn-prev`).disabled=!e.links?.prev,document.getElementById(`btn-next`).disabled=!e.links?.next,l=e.page}function h(e){if(!e)return``;let t=document.createElement(`div`);return t.innerText=e,t.innerHTML}async function g(e,t){w(e,`profiles`,`
    <button class="btn-back" onclick="window.history.back()">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M15 18l-6-6 6-6"/></svg>
      Back
    </button>
    <div id="detail-content" class="loading"><div class="spinner"></div> Loading profile...</div>
  `);try{let e=await(await r(`/api/profiles/${t.id}`)).json();if(e.status!==`success`)throw Error(e.message||`Profile not found`);let n=e.data;document.getElementById(`detail-content`).innerHTML=`
      <div class="profile-detail">
        <h2>${_(n.name)}</h2>
        
        <div class="detail-grid">
          <div class="detail-item">
            <span class="detail-label">ID</span>
            <span class="detail-value" style="font-family: monospace; font-size: 0.9rem;">${n.id}</span>
          </div>
          
          <div class="detail-item">
            <span class="detail-label">Gender</span>
            <span class="detail-value">
              <span class="badge badge-${n.gender}">${n.gender}</span>
              <span style="color: var(--text-muted); font-size: 0.85rem; margin-left: 8px;">
                ${(n.gender_probability*100).toFixed(0)}% confidence
              </span>
            </span>
          </div>

          <div class="detail-item">
            <span class="detail-label">Age</span>
            <span class="detail-value">${n.age} <span style="color: var(--text-muted); font-size: 0.85rem;">(${n.age_group})</span></span>
          </div>

          <div class="detail-item">
            <span class="detail-label">Country</span>
            <span class="detail-value">${_(n.country_name)} <span style="color: var(--text-muted); font-size: 0.85rem;">(${n.country_id})</span></span>
          </div>

          <div class="detail-item">
            <span class="detail-label">Country Confidence</span>
            <span class="detail-value">${(n.country_probability*100).toFixed(0)}%</span>
          </div>

          <div class="detail-item">
            <span class="detail-label">Created At</span>
            <span class="detail-value">${new Date(n.created_at).toLocaleString()}</span>
          </div>
        </div>
      </div>
    `}catch(e){document.getElementById(`detail-content`).innerHTML=`
      <div class="profile-detail">
        <p style="color: var(--danger);">${_(e.message)}</p>
      </div>
    `}}function _(e){if(!e)return``;let t=document.createElement(`div`);return t.innerText=e,t.innerHTML}async function v(e){w(e,`search`,`
    <div class="page-header">
      <h1>Natural Language Search</h1>
      <p>Search profiles using conversational queries (e.g., "young males from nigeria")</p>
    </div>

    <div class="search-container">
      <form id="search-form" class="search-box">
        <input type="text" id="search-input" class="search-input" placeholder="Enter your query..." autocomplete="off" />
        <button type="submit" class="btn-search">Search</button>
      </form>

      <div class="table-container" id="search-results" style="display: none;">
        <div class="table-toolbar">
          <h3 style="font-size: 1rem; font-weight: 600;">Results</h3>
        </div>
        <div style="overflow-x: auto;">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Gender</th>
                <th>Age</th>
                <th>Country</th>
              </tr>
            </thead>
            <tbody id="search-tbody"></tbody>
          </table>
        </div>
        <div class="pagination">
          <div class="pagination-info" id="search-pagination-info"></div>
        </div>
      </div>
    </div>
  `);let t=document.getElementById(`search-form`),n=document.getElementById(`search-input`),i=document.getElementById(`search-results`),a=document.getElementById(`search-tbody`);t.addEventListener(`submit`,async e=>{e.preventDefault();let t=n.value.trim();if(t){i.style.display=`block`,a.innerHTML=`<tr><td colspan="4" class="loading"><div class="spinner"></div> Searching...</td></tr>`,document.getElementById(`search-pagination-info`).textContent=``;try{let e=await(await r(`/api/profiles/search?q=${encodeURIComponent(t)}`)).json();if(e.status!==`success`)throw Error(e.message||`Search failed`);if(!e.data||e.data.length===0){a.innerHTML=`<tr><td colspan="4" style="padding: 32px; text-align: center; color: var(--text-muted);">No profiles matched your query</td></tr>`;return}a.innerHTML=e.data.map(e=>`
        <tr onclick="window.location.hash='#/profiles/${e.id}'">
          <td style="font-weight: 500;">${y(e.name)}</td>
          <td><span class="badge badge-${e.gender}">${e.gender}</span></td>
          <td>${e.age} <span style="color: var(--text-muted); font-size: 0.8rem;">(${e.age_group})</span></td>
          <td>${y(e.country_name)}</td>
        </tr>
      `).join(``),document.getElementById(`search-pagination-info`).textContent=`Found ${e.total} matches`}catch(e){a.innerHTML=`<tr><td colspan="4" style="padding: 24px; color: var(--danger);">${y(e.message)}</td></tr>`}}})}function y(e){if(!e)return``;let t=document.createElement(`div`);return t.innerText=e,t.innerHTML}async function b(e){w(e,`account`,`
    <div class="page-header">
      <h1>Account Settings</h1>
      <p>Manage your Insighta Labs+ profile</p>
    </div>
    <div id="account-content" class="loading"><div class="spinner"></div> Loading account...</div>
  `);try{let e=await(await r(`/auth/me`)).json();if(e.status!==`success`)throw Error(e.message||`Failed to load user info`);let t=e.data;document.getElementById(`account-content`).innerHTML=`
      <div class="account-card">
        <div class="account-header">
          <img src="${x(t.avatar_url)||`https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png`}" alt="Avatar" class="account-avatar" />
          <div>
            <div class="account-name">${x(t.username)}</div>
            <div class="account-role">${t.role}</div>
          </div>
        </div>
        
        <div class="detail-grid">
          <div class="detail-item">
            <span class="detail-label">Email</span>
            <span class="detail-value">${x(t.email)||`<em>Not provided</em>`}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Status</span>
            <span class="detail-value" style="color: ${t.is_active?`var(--success)`:`var(--danger)`};">
              ${t.is_active?`Active`:`Inactive`}
            </span>
          </div>
          <div class="detail-item" style="grid-column: 1 / -1;">
            <span class="detail-label">User ID</span>
            <span class="detail-value" style="font-family: monospace; font-size: 0.9rem;">${t.id}</span>
          </div>
        </div>

        <button id="btn-logout" class="btn-logout">Logout</button>
      </div>
    `,document.getElementById(`btn-logout`).addEventListener(`click`,async()=>{try{await r(`/auth/logout`,{method:`POST`}),window.location.hash=`#/login`}catch(e){alert(`Logout failed: `+e.message)}})}catch(e){document.getElementById(`account-content`).innerHTML=`
      <div class="account-card">
        <p style="color: var(--danger);">${x(e.message)}</p>
      </div>
    `}}function x(e){if(!e)return``;let t=document.createElement(`div`);return t.innerText=e,t.innerHTML}var S=document.getElementById(`app`);new e({"/login":()=>s(S),"/dashboard":()=>c(S),"/profiles":()=>f(S),"/profiles/:id":e=>g(S,e),"/search":()=>v(S),"/account":()=>b(S)}).start();var C={dashboard:`<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="9"></rect><rect x="14" y="3" width="7" height="5"></rect><rect x="14" y="12" width="7" height="9"></rect><rect x="3" y="16" width="7" height="5"></rect></svg>`,profiles:`<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>`,search:`<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>`,account:`<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>`};function w(e,t,n){document.getElementById(`main-content-area`)?(document.querySelectorAll(`.nav-link`).forEach(e=>{let n=e.getAttribute(`href`).replace(`#/`,``);e.classList.toggle(`active`,n===t||t===`profiles`&&n.startsWith(`profiles`))}),document.getElementById(`main-content-area`).innerHTML=n):e.innerHTML=`
      <div class="app-layout">
        <aside class="sidebar">
          <div class="sidebar-logo">Insighta <span>Labs+</span></div>
          <nav class="sidebar-nav">
            <a class="nav-link ${t===`dashboard`?`active`:``}" href="#/dashboard">
              ${C.dashboard} Dashboard
            </a>
            <a class="nav-link ${t===`profiles`?`active`:``}" href="#/profiles">
              ${C.profiles} Profiles
            </a>
            <a class="nav-link ${t===`search`?`active`:``}" href="#/search">
              ${C.search} Search
            </a>
          </nav>
          <div class="sidebar-footer">
            <a class="nav-link ${t===`account`?`active`:``}" href="#/account">
              ${C.account} Account
            </a>
          </div>
        </aside>
        <main class="main-content" id="main-content-area">
          ${n}
        </main>
      </div>
    `}