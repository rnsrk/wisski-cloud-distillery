(()=>{const e=e=>{const t=document.getElementsByTagName("h"+e);Array.from(t).forEach((e=>{void 0!==e.id&&""!==e.id&&e.appendChild((e=>{const t=document.createElement("a");return t.className="header-link",t.href="#"+e,t.innerHTML="#",t})(e.id))}))};new Array(6).fill(0).forEach(((t,n)=>e(n+1)));const t={date:e=>new Date(e.innerText).toISOString(),path:e=>{const t=e.innerText.split("/");return t[t.length-1]},pathbuilders:()=>{const e=window.pathbuilders,t=document.createElement("span");let o=!1;if(Object.keys(e).forEach((r=>{o=!0;const a=r+".xml",c=e[r];t.append(n(a,r,c,"application/xml")),t.append(document.createTextNode(" "))})),!o)return"(none)";const r=document.createElement("small");return r.append(document.createTextNode("(click to download)")),t.append(r),t}},n=(e,t,n,o)=>{const r=new Blob([n],{type:o??"text/plain"}),a=document.createElement("a");return a.target="_blank",a.download=e,a.href=URL.createObjectURL(r),a.append(document.createTextNode(t)),a};Object.keys(t).forEach((e=>{const n=t[e];document.querySelectorAll("code."+e).forEach((e=>{const t=n(e);if("string"==typeof t)return e.innerHTML="",void e.appendChild(document.createTextNode(t));e.parentNode.replaceChild(t,e)}))}))})();