!function(){function t(t){return t&&t.__esModule?t.default:t}var e="undefined"!=typeof globalThis?globalThis:"undefined"!=typeof self?self:"undefined"!=typeof window?window:"undefined"!=typeof global?global:{},n={},r={},i=e.parcelRequireafa4;function s(t,e){return new Promise(((n,r)=>{const i=new WebSocket(location.href.replace("http","ws"));i.onclose=n,i.onerror=r,i.onmessage=t=>e(t.data),i.onopen=()=>t(i)}))}null==i&&((i=function(t){if(t in n)return n[t].exports;if(t in r){var e=r[t];delete r[t];var i={id:t,exports:{}};return n[t]=i,e.call(i.exports,i,i.exports),i.exports}var s=new Error("Cannot find module '"+t+"'");throw s.code="MODULE_NOT_FOUND",s}).register=function(t,e){r[t]=e},e.parcelRequireafa4=i);const o=document.getElementsByClassName("remote-action");Array.from(o).forEach((t=>{const e=t.getAttribute("data-action"),n=t.getAttribute("data-force-reload"),r=t.getAttribute("data-param"),i=t.getAttribute("data-confirm-param"),o=i?document.querySelector(i):null,a=function(){var e,n;const r=null!==(n=parseInt(null!==(e=t.getAttribute("data-buffer"))&&void 0!==e?e:"",10))&&void 0!==n?n:0;return isFinite(r)&&r>0?r:0}(),u=function(){return!o||o.value===r};if(o){const e=()=>{u()?t.removeAttribute("disabled"):t.setAttribute("disabled","disabled")};o.addEventListener("change",e),e()}t.addEventListener("click",(function(t){if(t.preventDefault(),!u())return;const i=document.createElement("div");i.className="modal-terminal",document.body.append(i);const o=document.createElement("pre"),c=function(t,e,n){let r=null;const i=[],s=()=>{o.paintedFrames++,t.innerText=i.join("\n"),e.scrollTop=e.scrollHeight,r=null},o=(t,e)=>{if(i.push(t),0!==n&&i.length>n&&i.splice(0,i.length-n),null!==r&&(o.missedFrames++,window.cancelAnimationFrame(r)),e)return s();r=window.requestAnimationFrame(s)};return o.paintedFrames=0,o.missedFrames=0,o}(o,i,a);i.append(o);const d=document.createElement("button");d.className="pure-button pure-button-success",d.append("string"==typeof n?"Close & Reload":"Close"),d.addEventListener("click",(function(t){var e;if(t.preventDefault(),"string"==typeof n)return d.setAttribute("disabled","disabled"),o.innerHTML="Reloading page ...",void(""===n?location.reload():location.href=n);null===(e=i.parentNode)||void 0===e||e.removeChild(i)}));const l=window.onbeforeunload;window.onbeforeunload=()=>"A remote session is in progress. Are you sure you want to leave?";let f=!1;const h=function(){f||(f=!0,window.onbeforeunload=l,i.append(d))};c("Connecting ...",!0),s((t=>{c("Connected",!0),t.send(e),"string"==typeof r&&t.send(r)}),(t=>{c(t)})).then((()=>{c("Connection closed.",!0),h()})).catch((()=>{c("Connection errored.",!0),h()}))}))}));const a=document.getElementsByClassName("remote-link");Array.from(a).forEach((t=>{const e=t.getAttribute("data-action"),n=t.getAttribute("data-params"),r=null==n?void 0:n.split(" ");t.addEventListener("click",(function(t){t.preventDefault(),async function(t,e){return new Promise(((n,r)=>{let i="";var o=function(){const t=i.indexOf("\n");if(t<0)return void r("invalid buffer");if(!("true"===i.substring(0,t)))return void r(i);const e=JSON.parse(i.substring(t+1));n(e)};s((n=>{n.send(t),e&&e.forEach((t=>n.send(t)))}),(t=>{i+=t+"\n"})).then((()=>{o()})).catch((()=>{i="false\n",o()}))}))}(e,r).then((t=>{window.open(t)})).catch((t=>{console.error(t)}))}))}));var u={};u=function(){"use strict";var t=1e3,e=6e4,n=36e5,r="millisecond",i="second",s="minute",o="hour",a="day",u="week",c="month",d="quarter",l="year",f="date",h="Invalid Date",m=/^(\d{4})[-/]?(\d{1,2})?[-/]?(\d{0,2})[Tt\s]*(\d{1,2})?:?(\d{1,2})?:?(\d{1,2})?[.:]?(\d+)?$/,p=/\[([^\]]+)]|Y{1,4}|M{1,4}|D{1,2}|d{1,4}|H{1,2}|h{1,2}|a|A|m{1,2}|s{1,2}|Z{1,2}|SSS/g,$={name:"en",weekdays:"Sunday_Monday_Tuesday_Wednesday_Thursday_Friday_Saturday".split("_"),months:"January_February_March_April_May_June_July_August_September_October_November_December".split("_")},g=function(t,e,n){var r=String(t);return!r||r.length>=e?t:""+Array(e+1-r.length).join(n)+t},v={s:g,z:function(t){var e=-t.utcOffset(),n=Math.abs(e),r=Math.floor(n/60),i=n%60;return(e<=0?"+":"-")+g(r,2,"0")+":"+g(i,2,"0")},m:function t(e,n){if(e.date()<n.date())return-t(n,e);var r=12*(n.year()-e.year())+(n.month()-e.month()),i=e.clone().add(r,c),s=n-i<0,o=e.clone().add(r+(s?-1:1),c);return+(-(r+(n-i)/(s?i-o:o-i))||0)},a:function(t){return t<0?Math.ceil(t)||0:Math.floor(t)},p:function(t){return{M:c,y:l,w:u,d:a,D:f,h:o,m:s,s:i,ms:r,Q:d}[t]||String(t||"").toLowerCase().replace(/s$/,"")},u:function(t){return void 0===t}},y="en",w={};w[y]=$;var b=function(t){return t instanceof O},M=function t(e,n,r){var i;if(!e)return y;if("string"==typeof e){var s=e.toLowerCase();w[s]&&(i=s),n&&(w[s]=n,i=s);var o=e.split("-");if(!i&&o.length>1)return t(o[0])}else{var a=e.name;w[a]=e,i=a}return!r&&i&&(y=i),i||!r&&y},D=function(t,e){if(b(t))return t.clone();var n="object"==typeof e?e:{};return n.date=t,n.args=arguments,new O(n)},S=v;S.l=M,S.i=b,S.w=function(t,e){return D(t,{locale:e.$L,utc:e.$u,x:e.$x,$offset:e.$offset})};var O=function(){function $(t){this.$L=M(t.locale,null,!0),this.parse(t)}var g=$.prototype;return g.parse=function(t){this.$d=function(t){var e=t.date,n=t.utc;if(null===e)return new Date(NaN);if(S.u(e))return new Date;if(e instanceof Date)return new Date(e);if("string"==typeof e&&!/Z$/i.test(e)){var r=e.match(m);if(r){var i=r[2]-1||0,s=(r[7]||"0").substring(0,3);return n?new Date(Date.UTC(r[1],i,r[3]||1,r[4]||0,r[5]||0,r[6]||0,s)):new Date(r[1],i,r[3]||1,r[4]||0,r[5]||0,r[6]||0,s)}}return new Date(e)}(t),this.$x=t.x||{},this.init()},g.init=function(){var t=this.$d;this.$y=t.getFullYear(),this.$M=t.getMonth(),this.$D=t.getDate(),this.$W=t.getDay(),this.$H=t.getHours(),this.$m=t.getMinutes(),this.$s=t.getSeconds(),this.$ms=t.getMilliseconds()},g.$utils=function(){return S},g.isValid=function(){return!(this.$d.toString()===h)},g.isSame=function(t,e){var n=D(t);return this.startOf(e)<=n&&n<=this.endOf(e)},g.isAfter=function(t,e){return D(t)<this.startOf(e)},g.isBefore=function(t,e){return this.endOf(e)<D(t)},g.$g=function(t,e,n){return S.u(t)?this[e]:this.set(n,t)},g.unix=function(){return Math.floor(this.valueOf()/1e3)},g.valueOf=function(){return this.$d.getTime()},g.startOf=function(t,e){var n=this,r=!!S.u(e)||e,d=S.p(t),h=function(t,e){var i=S.w(n.$u?Date.UTC(n.$y,e,t):new Date(n.$y,e,t),n);return r?i:i.endOf(a)},m=function(t,e){return S.w(n.toDate()[t].apply(n.toDate("s"),(r?[0,0,0,0]:[23,59,59,999]).slice(e)),n)},p=this.$W,$=this.$M,g=this.$D,v="set"+(this.$u?"UTC":"");switch(d){case l:return r?h(1,0):h(31,11);case c:return r?h(1,$):h(0,$+1);case u:var y=this.$locale().weekStart||0,w=(p<y?p+7:p)-y;return h(r?g-w:g+(6-w),$);case a:case f:return m(v+"Hours",0);case o:return m(v+"Minutes",1);case s:return m(v+"Seconds",2);case i:return m(v+"Milliseconds",3);default:return this.clone()}},g.endOf=function(t){return this.startOf(t,!1)},g.$set=function(t,e){var n,u=S.p(t),d="set"+(this.$u?"UTC":""),h=(n={},n[a]=d+"Date",n[f]=d+"Date",n[c]=d+"Month",n[l]=d+"FullYear",n[o]=d+"Hours",n[s]=d+"Minutes",n[i]=d+"Seconds",n[r]=d+"Milliseconds",n)[u],m=u===a?this.$D+(e-this.$W):e;if(u===c||u===l){var p=this.clone().set(f,1);p.$d[h](m),p.init(),this.$d=p.set(f,Math.min(this.$D,p.daysInMonth())).$d}else h&&this.$d[h](m);return this.init(),this},g.set=function(t,e){return this.clone().$set(t,e)},g.get=function(t){return this[S.p(t)]()},g.add=function(r,d){var f,h=this;r=Number(r);var m=S.p(d),p=function(t){var e=D(h);return S.w(e.date(e.date()+Math.round(t*r)),h)};if(m===c)return this.set(c,this.$M+r);if(m===l)return this.set(l,this.$y+r);if(m===a)return p(1);if(m===u)return p(7);var $=(f={},f[s]=e,f[o]=n,f[i]=t,f)[m]||1,g=this.$d.getTime()+r*$;return S.w(g,this)},g.subtract=function(t,e){return this.add(-1*t,e)},g.format=function(t){var e=this,n=this.$locale();if(!this.isValid())return n.invalidDate||h;var r=t||"YYYY-MM-DDTHH:mm:ssZ",i=S.z(this),s=this.$H,o=this.$m,a=this.$M,u=n.weekdays,c=n.months,d=function(t,n,i,s){return t&&(t[n]||t(e,r))||i[n].slice(0,s)},l=function(t){return S.s(s%12||12,t,"0")},f=n.meridiem||function(t,e,n){var r=t<12?"AM":"PM";return n?r.toLowerCase():r},m={YY:String(this.$y).slice(-2),YYYY:this.$y,M:a+1,MM:S.s(a+1,2,"0"),MMM:d(n.monthsShort,a,c,3),MMMM:d(c,a),D:this.$D,DD:S.s(this.$D,2,"0"),d:String(this.$W),dd:d(n.weekdaysMin,this.$W,u,2),ddd:d(n.weekdaysShort,this.$W,u,3),dddd:u[this.$W],H:String(s),HH:S.s(s,2,"0"),h:l(1),hh:l(2),a:f(s,o,!0),A:f(s,o,!1),m:String(o),mm:S.s(o,2,"0"),s:String(this.$s),ss:S.s(this.$s,2,"0"),SSS:S.s(this.$ms,3,"0"),Z:i};return r.replace(p,(function(t,e){return e||m[t]||i.replace(":","")}))},g.utcOffset=function(){return 15*-Math.round(this.$d.getTimezoneOffset()/15)},g.diff=function(r,f,h){var m,p=S.p(f),$=D(r),g=($.utcOffset()-this.utcOffset())*e,v=this-$,y=S.m(this,$);return y=(m={},m[l]=y/12,m[c]=y,m[d]=y/3,m[u]=(v-g)/6048e5,m[a]=(v-g)/864e5,m[o]=v/n,m[s]=v/e,m[i]=v/t,m)[p]||v,h?y:S.a(y)},g.daysInMonth=function(){return this.endOf(c).$D},g.$locale=function(){return w[this.$L]},g.locale=function(t,e){if(!t)return this.$L;var n=this.clone(),r=M(t,e,!0);return r&&(n.$L=r),n},g.clone=function(){return S.w(this.$d,this)},g.toDate=function(){return new Date(this.valueOf())},g.toJSON=function(){return this.isValid()?this.toISOString():null},g.toISOString=function(){return this.$d.toISOString()},g.toString=function(){return this.$d.toUTCString()},$}(),A=O.prototype;return D.prototype=A,[["$ms",r],["$s",i],["$m",s],["$H",o],["$W",a],["$M",c],["$y",l],["$D",f]].forEach((function(t){A[t[1]]=function(e){return this.$g(e,t[0],t[1])}})),D.extend=function(t,e){return t.$i||(t(e,O,D),t.$i=!0),D},D.locale=M,D.isDayjs=b,D.unix=function(t){return D(1e3*t)},D.en=w[y],D.Ls=w,D.p={},D}();const c={date:e=>{const n=t(u)(e.innerText),r=n.format("YYYY-MM-DD HH:mm:ss ([UTC]Z)");if(0===n.unix()){const t=document.createElement("code");return t.style.color="gray",t.append(r),t}return r},path:t=>{const e=t.innerText.split("/");return e[e.length-1]},pathbuilder:t=>{var e;const n=(null!==(e=t.getAttribute("data-name"))&&void 0!==e?e:"pathbuilder")+".xml",[r,i]=d(n,t.innerText,"application/xml");r.className="pure-button";const s=n+" ("+i.size+" Bytes)";return r.append(s),r}},d=(t,e,n)=>{const r=new Blob([e],{type:null!=n?n:"text/plain"}),i=document.createElement("a");return i.target="_blank",i.download=t,i.href=URL.createObjectURL(r),[i,r]};Object.keys(c).forEach((t=>{const e=c[t];document.querySelectorAll("code."+t).forEach((t=>{const n=e(t);if("string"==typeof n)return t.innerHTML="",void t.appendChild(document.createTextNode(n));t.parentNode.replaceChild(n,t)}))})),i("kEAtK")}();