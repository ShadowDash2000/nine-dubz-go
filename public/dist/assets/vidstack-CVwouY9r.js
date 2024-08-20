import{bG as m,bH as c,bI as L,bJ as p,bK as S,bL as v,bM as g,bN as E,bO as u,bP as w,bQ as R,bR as y,bS as b,bT as H,bU as T,bV as I,bW as D,bX as M,bY as _}from"./index.js";const $=h=>_(h);class x{constructor(t,i){this.m=t,this.b=i,this.d=null,this.rb=null,this.sb={},this.tb=new Set}get instance(){return this.d}setup(t){const{streamType:i}=this.b.$state,e=p(i).includes("live"),r=p(i).includes("ll-");this.d=new t({lowLatencyMode:r,backBufferLength:r?4:e?8:void 0,renderTextTracksNatively:!1,...this.sb});const n=this.Pi.bind(this);for(const o of Object.values(t.Events))this.d.on(o,n);this.d.on(t.Events.ERROR,this.R.bind(this));for(const o of this.tb)o(this.d);this.b.player.dispatch("hls-instance",{detail:this.d}),this.d.attachMedia(this.m),this.d.on(t.Events.AUDIO_TRACK_SWITCHED,this.Ri.bind(this)),this.d.on(t.Events.LEVEL_SWITCHED,this.Si.bind(this)),this.d.on(t.Events.LEVEL_LOADED,this.Ti.bind(this)),this.d.on(t.Events.NON_NATIVE_TEXT_TRACKS_FOUND,this.Ui.bind(this)),this.d.on(t.Events.CUES_PARSED,this.Vi.bind(this)),this.b.qualities[v.Ja]=this.ke.bind(this),g(this.b.qualities,"change",this.le.bind(this)),g(this.b.audioTracks,"change",this.me.bind(this)),this.rb=E(this.ne.bind(this))}ba(t,i){return new u($(t),{detail:i})}ne(){if(!this.b.$state.live())return;const t=new w(this.oe.bind(this));return t.Ya(),t.aa.bind(t)}oe(){var t;this.b.$state.liveSyncPosition.set(((t=this.d)==null?void 0:t.liveSyncPosition)??1/0)}Pi(t,i){var e;(e=this.b.player)==null||e.dispatch(this.ba(t,i))}Ui(t,i){const e=this.ba(t,i);let r=-1;for(let n=0;n<i.tracks.length;n++){const o=i.tracks[n],s=o.subtitleTrack??o.closedCaptions,a=new R({id:`hls-${o.kind}-${n}`,src:s==null?void 0:s.url,label:o.label,language:s==null?void 0:s.lang,kind:o.kind,default:o.default});a[y.na]=2,a[y.ib]=()=>{a.mode==="showing"?(this.d.subtitleTrack=n,r=n):r===n&&(this.d.subtitleTrack=-1,r=-1)},this.b.textTracks.add(a,e)}}Vi(t,i){var o;const e=(o=this.d)==null?void 0:o.subtitleTrack,r=this.b.textTracks.getById(`hls-${i.type}-${e}`);if(!r)return;const n=this.ba(t,i);for(const s of i.cues)s.positionAlign="auto",r.addCue(s,n)}Ri(t,i){const e=this.b.audioTracks[i.id];if(e){const r=this.ba(t,i);this.b.audioTracks[b.fa](e,!0,r)}}Si(t,i){const e=this.b.qualities[i.level];if(e){const r=this.ba(t,i);this.b.qualities[b.fa](e,!0,r)}}Ti(t,i){var f;if(this.b.$state.canPlay())return;const{type:e,live:r,totalduration:n,targetduration:o}=i.details,s=this.ba(t,i);this.b.delegate.c("stream-type-change",r?e==="EVENT"&&Number.isFinite(n)&&o>=10?"live:dvr":"live":"on-demand",s),this.b.delegate.c("duration-change",n,s);const a=this.d.media;this.d.currentLevel===-1&&this.b.qualities[v.Xa](!0,s);for(const d of this.d.audioTracks){const l={id:d.id.toString(),label:d.name,language:d.lang||"",kind:"main"};this.b.audioTracks[b.ea](l,s)}for(const d of this.d.levels){const l={id:((f=d.id)==null?void 0:f.toString())??d.height+"p",width:d.width,height:d.height,codec:d.codecSet,bitrate:d.bitrate};this.b.qualities[b.ea](l,s)}a.dispatchEvent(new u("canplay",{trigger:s}))}R(t,i){var e;if(i.fatal)switch(i.type){case"mediaError":(e=this.d)==null||e.recoverMediaError();break;default:this.rc(i.error);break}}rc(t){this.b.delegate.c("error",{message:t.message,code:1,error:t})}ke(){this.d&&(this.d.currentLevel=-1)}le(){const{qualities:t}=this.b;!this.d||t.auto||(this.d[t.switch+"Level"]=t.selectedIndex,H&&(this.m.currentTime=this.m.currentTime))}me(){const{audioTracks:t}=this.b;this.d&&this.d.audioTrack!==t.selectedIndex&&(this.d.audioTrack=t.selectedIndex)}Wi(t){var i;c(t.src)&&((i=this.d)==null||i.loadSource(t.src))}Xi(){var t,i;(t=this.d)==null||t.destroy(),this.d=null,(i=this.rb)==null||i.call(this),this.rb=null}}class C{constructor(t,i,e){this.M=t,this.b=i,this.Ma=e,this.re()}async re(){const t={onLoadStart:this.Na.bind(this),onLoaded:this.ub.bind(this),onLoadError:this.se.bind(this)};let i=await O(this.M,t);if(T(i)&&!c(this.M)&&(i=await N(this.M,t)),!i)return null;if(!i.isSupported()){const e="[vidstack] `hls.js` is not supported in this environment";return this.b.player.dispatch(new u("hls-unsupported")),this.b.delegate.c("error",{message:e,code:4}),null}return i}Na(){this.b.player.dispatch(new u("hls-lib-load-start"))}ub(t){this.b.player.dispatch(new u("hls-lib-loaded",{detail:t})),this.Ma(t)}se(t){const i=I(t);this.b.player.dispatch(new u("hls-lib-load-error",{detail:i})),this.b.delegate.c("error",{message:i.message,code:4,error:i})}}async function N(h,t={}){var i,e,r,n,o;if(!T(h)){if((i=t.onLoadStart)==null||i.call(t),h.prototype&&h.prototype!==Function)return(e=t.onLoaded)==null||e.call(t,h),h;try{const s=(r=await h())==null?void 0:r.default;if(s&&s.isSupported)(n=t.onLoaded)==null||n.call(t,s);else throw Error("");return s}catch(s){(o=t.onLoadError)==null||o.call(t,s)}}}async function O(h,t={}){var i,e,r;if(c(h)){(i=t.onLoadStart)==null||i.call(t);try{if(await D(h),!M(window.Hls))throw Error("");const n=window.Hls;return(e=t.onLoaded)==null||e.call(t,n),n}catch(n){(r=t.onLoadError)==null||r.call(t,n)}}}const P="https://cdn.jsdelivr.net";class V extends m{constructor(){super(...arguments),this.$$PROVIDER_TYPE="HLS",this.sc=null,this.e=new x(this.video,this.b),this.pa=`${P}/npm/hls.js@^1.5.0/dist/hls.min.js`}get ctor(){return this.sc}get instance(){return this.e.instance}get type(){return"hls"}get canLiveSync(){return!0}get config(){return this.e.sb}set config(t){this.e.sb=t}get library(){return this.pa}set library(t){this.pa=t}preconnect(){c(this.pa)&&L(this.pa)}setup(){super.setup(),new C(this.pa,this.b,t=>{this.sc=t,this.e.setup(t),this.b.delegate.c("provider-setup",this);const i=p(this.b.$state.source);i&&this.loadSource(i)})}async loadSource(t,i){if(!c(t.src)){this.pc();return}this.a.preload=i||"",this.he(t,"application/x-mpegurl"),this.e.Wi(t),this.L=t}onInstance(t){const i=this.e.instance;return i&&t(i),this.e.tb.add(t),()=>this.e.tb.delete(t)}destroy(){this.e.Xi()}}V.supported=S();export{V as HLSProvider};
