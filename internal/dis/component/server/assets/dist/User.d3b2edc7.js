var e,t,r,n;e="undefined"!=typeof globalThis?globalThis:"undefined"!=typeof self?self:"undefined"!=typeof window?window:"undefined"!=typeof global?global:{},t={},r={},null==(n=e.parcelRequireafa4)&&((n=function(e){if(e in t)return t[e].exports;if(e in r){var n=r[e];delete r[e];var o={id:e,exports:{}};return t[e]=o,n.call(o.exports,o,o.exports),o.exports}var i=Error("Cannot find module '"+e+"'");throw i.code="MODULE_NOT_FOUND",i}).register=function(e,t){r[e]=t},e.parcelRequireafa4=n),n.register("eejBl",function(e,t){n("5J9YX"),n("lNUmu")}),n.register("5J9YX",function(e,t){var r=n("aPvYT");document.querySelectorAll(".copy").forEach(e=>{e.addEventListener("click",()=>{navigator.clipboard&&(0,r.discard)(navigator.clipboard.writeText(e.innerText))})})}),n.register("aPvYT",function(e,t){Object.defineProperty(e.exports,"discard",{get:function(){return n},set:void 0,enumerable:!0,configurable:!0});let r=console.error.bind(console);function n(e){e.then(()=>{}).catch(r)}}),n.register("lNUmu",function(e,t){document.querySelectorAll("span").forEach(e=>{e.hasAttribute("data-reveal")&&function(e,t){let r=e.getAttribute("data-reveal")??"(no content)",n=!0,o=()=>{n=!0,e.innerText="(click to reveal)"};o();let i=()=>{n=!1;let t=document.createElement("code");t.append(r),t.addEventListener("click",e=>{e.preventDefault(),navigator.clipboard&&navigator.clipboard.writeText(r).then(()=>{}).catch(console.error.bind(console))}),t.style.userSelect="all",e.innerHTML="",e.append(t)};e.addEventListener("click",e=>{e.preventDefault(),n&&(i(),setTimeout(o,1e4))})}(e,0)})}),n("eejBl");
//# sourceMappingURL=User.d3b2edc7.js.map