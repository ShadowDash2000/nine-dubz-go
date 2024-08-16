import{ce as m,a as x,c as g,F as r,j as e,L as i,I as l,T as u,d,B as j}from"./index.js";import{u as f,a as p,b as h,c as y,d as b,e as w,L as I,f as N,R}from"./configForm.js";const c=m(n=>({message:"",login:async o=>{try{const t=await(await fetch("/api/authorize/inner/login",{method:"POST",body:JSON.stringify(o)})).json();return n({message:t==null?void 0:t.message}),t}catch(a){console.log(a.message)}}})),v=()=>{const n=c(s=>s.login),o=c(s=>s.message),a=x();return{signIn:g({mutationKey:["login"],mutationFn:s=>n(s),onSuccess:s=>{if((s==null?void 0:s.status)!=="error")return a("/")}}),message:o}},S=()=>{const[n]=r.useForm(),o=f(),{signIn:a,message:t}=v();return e.jsxs(r,{className:"min-w-80",onFinish:a.mutate,form:n,layout:"vertical",children:[e.jsx(r.Item,{children:e.jsxs(i,{to:o,className:`flex w-full gap-2 justify-center rounded-md border\r
                            border-gray-300 bg-white py-3 px-4 text-sm font-medium text-black shadow-sm hover:bg-gray-50`,children:[e.jsx(p,{style:{fontSize:"20px"}}),"Sign In with Google"]})}),e.jsx(r.Item,{children:e.jsxs("div",{className:"relative flex items-center",children:[e.jsx("div",{className:"flex-grow border-t border-gray-200"}),e.jsx("span",{className:"mx-4 text-xs text-gray-200 text-center uppercase whitespace-nowrap",children:"Or"}),e.jsx("div",{className:"flex-grow border-t border-gray-200"})]})}),e.jsx(r.Item,{className:"text-gray-200",label:"Email",name:"email",rules:[h],children:e.jsx(l,{size:"large",prefix:e.jsx(y,{})})}),e.jsx(r.Item,{className:"text-gray-200",label:"Password",name:"password",rules:[b],children:e.jsx(l,{size:"large",prefix:e.jsx(w,{}),suffix:e.jsx(u,{placement:"right",title:e.jsx(I,{listRiles:N,field:"Пароль"}),children:e.jsx(R,{})})})}),t&&e.jsx("small",{className:"text-red-500 block mt-2",children:t}),e.jsx(r.Item,{children:e.jsxs(d,{align:"center",justify:"space-between",children:[e.jsx(j,{type:"default",htmlType:"submit",children:"Войти"}),e.jsx(i,{className:"text-gray-200",to:"/signup",children:"Еще нет аккаунта?"})]})})]})};export{S as default};
