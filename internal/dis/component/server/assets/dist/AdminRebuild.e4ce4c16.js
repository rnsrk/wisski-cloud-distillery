!function(){var e="undefined"!=typeof globalThis?globalThis:"undefined"!=typeof self?self:"undefined"!=typeof window?window:"undefined"!=typeof global?global:{},t={},n={},o=e.parcelRequireafa4;null==o&&((o=function(e){if(e in t)return t[e].exports;if(e in n){var o=n[e];delete n[e];var r={id:e,exports:{}};return t[e]=r,o.call(r.exports,r,r.exports),r.exports}var i=Error("Cannot find module '"+e+"'");throw i.code="MODULE_NOT_FOUND",i}).register=function(e,t){n[e]=t},e.parcelRequireafa4=o),o("dK5Bi");var r=o("8vh0V");async function i(e,t){return await new Promise((n,o)=>{(0,r.createModal)("rebuild",[e,JSON.stringify(t)],{bufferSize:0,onClose:(t,r)=>{if(!t){o(Error(r??"unspecified error"));return}n(e)}})})}let l=document.getElementById("system"),d=document.getElementById("slug"),a=document.getElementById("php"),u=document.getElementById("opcacheDevelopment"),c=document.getElementById("contentsecuritypolicy");l.addEventListener("submit",e=>{e.preventDefault(),i(d.value,{PHP:a.value,OpCacheDevelopment:u.checked,ContentSecurityPolicy:c.value}).then(e=>{location.href="/admin/instance/"+e}).catch(e=>{console.error(e),location.reload()})}),l.querySelector("fieldset")?.removeAttribute("disabled")}();