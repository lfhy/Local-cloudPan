import{an as V,c3 as M,aR as A,a5 as F,Q as p,bM as N,o as a,c,bA as O,ay as P,aN as y,b6 as S,b as r,aA as $,aH as d,R as _,ax as C,aF as i,aB as u,c6 as U,c7 as j,aI as h,aD as q}from"./index-B5pPcsYu.js";import{g as G,E as H,a as Q}from"./index-DU-i08af.js";const T=["onContextmenu"],J=["onUpdate:modelValue"],K=["onClick"],W=["src"],X=["value"],Y={key:3,class:"filename"},Z={class:"describe"},ee={key:0},te={key:1},ce=V({name:"GridView",__name:"index",props:{fileList:{}},emits:["fileClick","fileChoosed","finish","cancel"],setup(b,{expose:x,emit:E}){const m=b,l=E,w=M(),s=A(m.fileList.reduce((e,n)=>(e[n.id]=!1,e),{})),f=F(null),g=e=>l("fileClick",m.fileList[e]),z=(e,n)=>{e.preventDefault(),s[n]=!0},B=async e=>{l("finish",e,f.value[0].value),s[e.id]=!1},L=e=>{l("cancel",e),s[e.id]=!1};return p(s,()=>l("fileChoosed",s)),p(()=>m.fileList,e=>{e.forEach(n=>{s[n.id]===void 0&&(s[n.id]=!1)})}),p(()=>w.query.path,()=>{Object.keys(s).forEach(e=>{delete s[e]})}),x({clearSelection:()=>Object.keys(s).map(e=>s[e]=!1)}),(e,n)=>{const D=H,k=q,v=Q,I=N("focus");return a(!0),c(C,null,O(e.fileList,(t,R)=>(a(),c("div",{key:t.id,class:P(["file-item",{checked:s[t.id]}]),onContextmenu:o=>z(o,t.id)},[y(r("input",{class:"file-checkbox",type:"checkbox","onUpdate:modelValue":o=>s[t.id]=o,name:"selected"},null,8,J),[[S,s[t.id]]]),r("div",{onClick:o=>g(R)},[t.fileType==="picture"?(a(),$(D,{key:0,src:t.thumbnailPath||t.filePath,lazy:!0,"preview-src-list":[t.filePath],class:"picture",fit:"cover","preview-teleported":"",onClick:d(()=>{},["stop"])},null,8,["src","preview-src-list"])):(a(),c("img",{key:1,src:_(G)(t),class:"picture"},null,8,W)),t.isRename?(a(),c(C,{key:2},[y(r("input",{class:"rename-ipt",type:"text",name:"filename",ref_for:!0,ref_key:"renameIpt",ref:f,value:t.name,onClick:d(()=>{},["stop"])},null,8,X),[[I]]),i(v,{type:"primary",size:"small",class:"rename-btn",onClick:d(o=>B(t),["stop"])},{default:u(()=>[i(k,{size:"12"},{default:u(()=>[i(_(U))]),_:1})]),_:2},1032,["onClick"]),i(v,{type:"primary",size:"small",class:"rename-btn",onClick:d(o=>L(t),["stop"])},{default:u(()=>[i(k,{size:"12"},{default:u(()=>[i(_(j))]),_:1})]),_:2},1032,["onClick"])],64)):(a(),c("p",Y,h(t.name),1)),r("div",Z,[t.isDir?(a(),c("p",ee,h(t.modified),1)):(a(),c("p",te,h(t.size),1))])],8,K)],42,T))),128)}}});export{ce as default};
