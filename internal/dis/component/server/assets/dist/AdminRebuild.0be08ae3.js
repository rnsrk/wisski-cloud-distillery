!function(){var e="undefined"!=typeof globalThis?globalThis:"undefined"!=typeof self?self:"undefined"!=typeof window?window:"undefined"!=typeof global?global:{},t={},n={},r=e.parcelRequireafa4;null==r&&((r=function(e){if(e in t)return t[e].exports;if(e in n){var r=n[e];delete n[e];var o={id:e,exports:{}};return t[e]=o,r.call(o.exports,o,o.exports),o.exports}var i=Error("Cannot find module '"+e+"'");throw i.code="MODULE_NOT_FOUND",i}).register=function(e,t){n[e]=t},e.parcelRequireafa4=r),r("dK5Bi");var o=r("8vh0V");async function i(e,t){return await new Promise((n,r)=>{(0,o.createModal)("rebuild",[e,JSON.stringify(t)],{bufferSize:0,onClose:(t,o)=>{if(!t){r(Error(o??"unspecified error"));return}n(e)}})})}let l=document.getElementById("system"),d=document.getElementById("slug"),a=document.getElementById("php"),c=document.getElementById("opcacheDevelopment"),u=document.getElementById("contentsecuritypolicy"),f=document.getElementById("iipserver");l.addEventListener("submit",e=>{e.preventDefault(),i(d.value,{PHP:a.value,IIPServer:f.checked,OpCacheDevelopment:c.checked,ContentSecurityPolicy:u.value}).then(e=>{location.href="/admin/instance/"+e}).catch(e=>{console.error(e),location.reload()})}),l.querySelector("fieldset")?.removeAttribute("disabled")}();