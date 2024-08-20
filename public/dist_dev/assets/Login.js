import{ca as m,a as x,c as g,F as r,j as e,L as l,I as i,d as u,e as d,B as h}from"./index.js";import{u as f,R as j,a as p,b as y,c as w,d as b}from"./configForm.js";import{F as I}from"./index2.js";const c=m(o=>({message:"",login:async n=>{try{const t=await(await fetch("http://localhost:25565/api/authorize/inner/login",{method:"POST",body:JSON.stringify(n)})).json();return o({message:t==null?void 0:t.message}),t}catch(a){console.log(a.message)}}})),N=()=>{const o=c(s=>s.login),n=c(s=>s.message),a=x();return{signIn:g({mutationKey:["login"],mutationFn:s=>o(s),onSuccess:s=>{if((s==null?void 0:s.status)!=="error")return a("/")}}),message:n}},L=()=>{const[o]=r.useForm(),n=f(),{signIn:a,message:t}=N();return e.jsxs(r,{className:"min-w-80",onFinish:a.mutate,form:o,layout:"vertical",children:[e.jsx(r.Item,{children:e.jsxs(l,{to:n,className:`flex w-full gap-2 justify-center rounded-md border\r
                            border-gray-300 bg-white py-3 px-4 text-sm font-medium text-black shadow-sm hover:bg-gray-50`,children:[e.jsx(j,{style:{fontSize:"20px"}}),"Sign In with Google"]})}),e.jsx(r.Item,{children:e.jsxs("div",{className:"relative flex items-center",children:[e.jsx("div",{className:"flex-grow border-t border-gray-200"}),e.jsx("span",{className:"mx-4 text-xs text-gray-200 text-center uppercase whitespace-nowrap",children:"Or"}),e.jsx("div",{className:"flex-grow border-t border-gray-200"})]})}),e.jsx(r.Item,{className:"text-gray-200",label:"Email",name:"email",rules:[p],children:e.jsx(i,{size:"large",prefix:e.jsx(y,{})})}),e.jsx(r.Item,{className:"text-gray-200",label:"Password",name:"password",rules:[w],children:e.jsx(i.Password,{autocomplete:"off",size:"large",prefix:e.jsx(b,{}),iconRender:s=>s?e.jsx(u,{style:{color:"white"}}):e.jsx(d,{style:{color:"white"}})})}),t&&e.jsx("small",{className:"text-red-500 block mt-2",children:t}),e.jsx(r.Item,{children:e.jsxs(I,{align:"center",justify:"space-between",children:[e.jsx(h,{type:"default",htmlType:"submit",children:"Войти"}),e.jsx(l,{className:"text-gray-200",to:"/signup",children:"Еще нет аккаунта?"})]})})]})};export{L as default};
