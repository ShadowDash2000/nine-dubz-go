import{c0 as n,b$ as a,bM as o,bN as h,bJ as u,c5 as b,c6 as d,bH as l}from"./index.js";function p(e,t=3e3){const i=n();return setTimeout(()=>{const s=e();s&&i.reject(s)},t),i}class g{constructor(t){this.Mb=t,this.tc=a(""),this.referrerPolicy=null,t.setAttribute("frameBorder","0"),t.setAttribute("aria-hidden","true"),t.setAttribute("allow","autoplay; fullscreen; encrypted-media; picture-in-picture; accelerometer; gyroscope"),this.referrerPolicy!==null&&t.setAttribute("referrerpolicy",this.referrerPolicy)}get iframe(){return this.Mb}setup(){o(window,"message",this.Yi.bind(this)),o(this.Mb,"load",this.hd.bind(this)),h(this.Nb.bind(this))}Nb(){const t=this.tc();if(!t.length){this.Mb.setAttribute("src","");return}const i=u(()=>this.ng());this.Mb.setAttribute("src",b(t,i))}te(t,i){var s;d||(s=this.Mb.contentWindow)==null||s.postMessage(JSON.stringify(t),i??"*")}Yi(t){var c;const i=this.Ob();if((t.source===null||t.source===((c=this.Mb)==null?void 0:c.contentWindow))&&(!l(i)||i===t.origin)){try{const r=JSON.parse(t.data);r&&this.ue(r,t);return}catch{}t.data&&this.ue(t.data,t)}}}export{g as E,p as t};
