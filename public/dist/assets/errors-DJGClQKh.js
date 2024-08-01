import{P as t,c as a}from"./prod-CDLwOxs6.js";import"./index-C-cBWuF2.js";const s={p(){return new t({code:a.BadSignature,reason:"missing WEBVTT file header",line:1})},q(n,e){return new t({code:a.BadTimestamp,reason:`cue start timestamp \`${n}\` is invalid on line ${e}`,line:e})},r(n,e){return new t({code:a.BadTimestamp,reason:`cue end timestamp \`${n}\` is invalid on line ${e}`,line:e})},s(n,e,r){return new t({code:a.BadTimestamp,reason:`cue end timestamp \`${e}\` is greater than start \`${n}\` on line ${r}`,line:r})},w(n,e,r){return new t({code:a.BadSettingValue,reason:`invalid value for cue setting \`${n}\` on line ${r} (value: ${e})`,line:r})},v(n,e,r){return new t({code:a.UnknownSetting,reason:`unknown cue setting \`${n}\` on line ${r} (value: ${e})`,line:r})},u(n,e,r){return new t({code:a.BadSettingValue,reason:`invalid value for region setting \`${n}\` on line ${r} (value: ${e})`,line:r})},t(n,e,r){return new t({code:a.UnknownSetting,reason:`unknown region setting \`${n}\` on line ${r} (value: ${e})`,line:r})},N(n,e){return new t({code:a.BadFormat,reason:`format missing for \`${n}\` block on line ${e}`,line:e})}};export{s as ParseErrorBuilder};
