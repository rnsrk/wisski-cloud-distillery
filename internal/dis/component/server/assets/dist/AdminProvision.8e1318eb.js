!function(){var e="undefined"!=typeof globalThis?globalThis:"undefined"!=typeof self?self:"undefined"!=typeof window?window:"undefined"!=typeof global?global:{},n={},t={},o=e.parcelRequireafa4;null==o&&((o=function(e){if(e in n)return n[e].exports;if(e in t){var o=t[e];delete t[e];var r={id:e,exports:{}};return n[e]=r,o.call(r.exports,r,r.exports),r.exports}var l=new Error("Cannot find module '"+e+"'");throw l.code="MODULE_NOT_FOUND",l}).register=function(e,n){t[e]=n},e.parcelRequireafa4=o),o("dK5Bi");var r,l=o("8vh0V");async function i(e){return await new Promise(((n,t)=>{(0,l.createModal)("provision",[JSON.stringify(e)],{bufferSize:0,onClose:(o,r)=>{o?n(e.Slug):t(new Error(null!=r?r:"unspecified error"))}})}))}const d=document.getElementById("system"),a=document.getElementById("slug"),u=document.getElementById("php"),c=document.getElementById("opcacheDevelopment"),s=document.getElementById("contentsecuritypolicy");d.addEventListener("submit",(e=>{e.preventDefault(),i({Slug:a.value,System:{PHP:u.value,OpCacheDevelopment:c.checked,ContentSecurityPolicy:s.value}}).then((e=>{location.href="/admin/instance/"+e})).catch((e=>{console.error(e),location.reload()}))})),null===(r=d.querySelector("fieldset"))||void 0===r||r.removeAttribute("disabled")}();