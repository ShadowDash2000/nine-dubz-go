import{q as c,o as a,l as n,e as h,b as u,w as d,x as l,i as b}from"./index-C-cBWuF2.js";function p(r,t=3e3){const i=c();return setTimeout(()=>{const s=r();s&&i.reject(s)},t),i}class g{constructor(t){this.Mb=t,this.tc=a(""),this.referrerPolicy=null,t.setAttribute("frameBorder","0"),t.setAttribute("aria-hidden","true"),t.setAttribute("allow","autoplay; fullscreen; encrypted-media; picture-in-picture; accelerometer; gyroscope"),this.referrerPolicy!==null&&t.setAttribute("referrerpolicy",this.referrerPolicy)}get iframe(){return this.Mb}setup(){n(window,"message",this.Yi.bind(this)),n(this.Mb,"load",this.hd.bind(this)),h(this.Nb.bind(this))}Nb(){const t=this.tc();if(!t.length){this.Mb.setAttribute("src","");return}const i=u(()=>this.ng());this.Mb.setAttribute("src",d(t,i))}te(t,i){var s;l||(s=this.Mb.contentWindow)==null||s.postMessage(JSON.stringify(t),i??"*")}Yi(t){var o;const i=this.Ob();if((t.source===null||t.source===((o=this.Mb)==null?void 0:o.contentWindow))&&(!b(i)||i===t.origin)){try{const e=JSON.parse(t.data);e&&this.ue(e,t);return}catch{}t.data&&this.ue(t.data,t)}}}export{g as E,p as t};
