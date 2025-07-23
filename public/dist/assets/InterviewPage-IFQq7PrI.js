import{_ as ye,r as t,o as fe,a as be,c as v,b as o,n as M,d as Ie,t as I,e as ke,f as Be,F as Ae,g as De,h as c,u as Te,i as d}from"./index-AjIqghy1.js";const xe="/assets/%E5%A3%B0%E9%9F%B3%E5%85%B3%E9%97%AD-DpSzKatq.png",$e="/assets/%E5%BC%80%E5%90%AF%E5%A3%B0%E9%9F%B3-qtu0_L9e.png",Ce={class:"interview-container"},Re={class:"header"},Se={class:"info"},ze={class:"header-buttons"},Ee={class:"main-content"},Ue={class:"left"},Ve={class:"video-panel"},qe={class:"video-main"},Le={class:"recording-panel"},Me={key:0,src:xe,alt:"开始录音",class:"record-img"},Pe={key:1,src:$e,alt:"停止录音",class:"record-img"},Fe={class:"side-panel"},je={class:"ai-message-box"},Ne={class:"interaction-panel"},Oe={class:"conversation-area"},We={class:"message-content"},He={key:0,class:"audio-player"},Ke=["src"],Ge={key:1,class:"message-text"},Je={class:"message-time"},Qe={class:"ai-status"},Xe={key:0,class:"loading-overlay"},Ye={__name:"InterviewPage",setup(Ze){const X=Te(),x=t(null),k=t(null),$=t(null),s=t(localStorage.getItem("token")||""),i=t(localStorage.getItem("username")||""),f=t(!1),B=t(!1),p=t([]),g=t(null),Y=t(!1),u=t(""),A=t([]),m=t(null),h=t(null),_=t(!1),D=t(null),S=t(null),z=t(null),w=t(null),b=t([]),P=t(!1),C=t(null),E=t(null),T=t(""),U=t(""),F=t(null),R=t(""),V=t(""),q=t(!1),j=t(!1),L=t(!1);t(null);const y=t(null),N=t(!1),Z=async()=>{console.log("当前用户名:",i.value),console.log("当前token:",s.value);try{const e={username:i.value,prefer_role:0};console.log("发送的请求数据:",e);const a=await c.post("/api/conversation/new",e,{headers:{Authorization:s?`Bearer ${s.value}`:""}});k.value=a.data.data.conversation_id,console.log("成功创建对话，ID:",k.value)}catch(e){console.error("创建会话失败:",e)}},ee=async()=>{const e={username:i.value};try{const a=await c.post("api/interview/new",e,{headers:{Authorization:s.value?`Bearer ${s.value}`:""}});console.log("获取面试ID成功:",a.data.data),$.value=a.data.data,localStorage.setItem("interview_id",$.value)}catch(a){console.error("获取面试ID失败:",a)}},ae=async()=>{const e={username:i.value,interview_id:$.value,interview_conversation:R.value,video_url:E.value,emotion_conversation:V.value};try{const a=await c.post("api/interview/update",e,{headers:{Authorization:s.value?`Bearer ${s.value}`:""}});console.log("更新面试数据成功:",a.data)}catch(a){console.error("更新面试数据失败:",a)}},te=async()=>{console.log("面试结束，进入结果界面"),q.value=!0;try{K(),W(),await he(),await me(),await ae(),X.push(`/result/${$.value}`)}catch(e){console.error("跳转结果页失败:",e),q.value=!1,alert("生成报告失败，请重试")}},O=(e,a)=>{U.value="";let n=0;const l=setInterval(()=>{U.value+=e[n],n++,n>=e.length&&(clearInterval(l),a&&a())},80)},oe=async()=>{if(B.value){alert("面试已开始");return}await H(),console.log("是否开启白噪音：",L.value),L.value&&se(),ee(),re(),B.value=!0,f.value=!0,console.log("开始面试"),await Z();const e="你好，请你先介绍一下自己";T.value="",O(e,()=>{f.value=!1,T.value=e}),setTimeout(()=>{g.value&&(g.value.scrollTop=g.value.scrollHeight,p.value.push({type:"ai",text:e,time:new Date().toLocaleTimeString()}))},1e3)},se=()=>{y.value||(y.value=new Audio("../assets/audio/noisy.mp3"),y.value.loop=!0),y.value.play().then(()=>{N.value=!0,console.log("白噪音开始播放")}).catch(e=>{console.error("白噪音播放失败:",e)})},W=()=>{y.value&&(y.value.pause(),y.value.currentTime=0,N.value=!1,console.log("白噪音已停止"))},H=async()=>{try{const e=await navigator.mediaDevices.getUserMedia({video:!0,audio:!0});z.value=e,h.value&&(h.value.srcObject=e)}catch(e){console.error("无法访问摄像头:",e),alert("无法访问摄像头，请检查权限。")}},ne=()=>{z.value&&(b.value=[],w.value=new MediaRecorder(z.value,{mimeType:"video/webm; codecs=vp9"}),w.value.ondataavailable=e=>{e.data.size>0&&b.value.push(e.data)},w.value.onstop=()=>{const e=new Blob(b.value,{type:"video/webm"});C.value=URL.createObjectURL(e),console.log("Video recording stopped",C.value)},w.value.start(),P.value=!0,console.log("Video recording started"))},le=()=>{w.value&&w.value.state==="recording"&&(w.value.stop(),P.value=!1,console.log("Video recording stopped"))},re=()=>{console.log("开始录制30s视频用于情感分析"),ne(),F.value=setTimeout(async()=>{await le(),await new Promise(e=>{const a=setInterval(()=>{C.value&&(clearInterval(a),e())},100)}),ie(),F.value=null},1e4)},ie=async()=>{if(console.log("开始保存视频片段"),!C.value){console.error("视频URL不可用");return}try{if(!b.value||b.value.length===0){console.error("没有可用的视频片段");return}const e=new FormData,a=new Blob(b.value,{type:"video/webm"});console.log("videotype",a.type);const n={file:a};e.append("video",a,"30-second-clip.webm");const l=await c.post("/api/upload/video",n,{headers:{Authorization:s.value?`Bearer ${s.value}`:"","Content-Type":"multipart/form-data"}});console.log("视频上传成功:",l.data),E.value=l.data.data}catch(e){console.error("视频上传失败:",e)}},K=()=>{h.value&&h.value.srcObject&&(h.value.srcObject.getTracks().forEach(e=>e.stop()),h.value.srcObject=null)},ce=e=>new Promise((a,n)=>{const l=new FileReader;l.onloadend=()=>{a(l.result)},l.onerror=n,l.readAsDataURL(e)}),ue=async()=>{if(!u.value){console.error("没有音频数据可上传");return}try{const e=await c.post("api/upload/audio64",{base64:u.value},{headers:{Authorization:s.value?`Bearer ${s.value}`:""}});S.value=e.data,console.log("音频上传成功:",S.value.data)}catch(e){console.error("Error uploading audio:",e)}},ve=async()=>{if(B.value)try{D.value=await navigator.mediaDevices.getUserMedia({audio:!0}),A.value=[],m.value=new MediaRecorder(D.value,{mimeType:"audio/webm"}),m.value.ondataavailable=e=>{e.data.size>0&&A.value.push(e.data)},m.value.onstop=async()=>{const e=new Blob(A.value,{type:"audio/webm"});u.value=await ce(e),console.log("音频转换为Base64:",u.value),A.value=[]},m.value.start(),_.value=!0}catch(e){_.value=!1,console.error("无法访问麦克风或不支持录音:",e)}else alert("请先开始面试")},G=()=>{m.value&&m.value.state!=="inactive"&&m.value.stop(),D.value&&(D.value.getTracks().forEach(e=>e.stop()),D.value=null),_.value=!1},de=async()=>{if(_.value){if(G(),console.log("停止录音"),await new Promise(e=>{const a=setInterval(()=>{!_.value&&u.value&&(clearInterval(a),e())},100)}),u.value){p.value=[...p.value,{type:"base64",data:u.value,time:new Date().toLocaleTimeString()}],Y.value=!0;try{await ue(),await ge()}catch(e){console.error("处理音频时出错:",e)}u.value=""}}else console.log("开始录音"),A.value=[],await ve()},pe=async()=>{try{const e=await c.get("api/user/get",{headers:{Authorization:`Bearer ${s.value}`},params:{username:i.value}});x.value=e.data,L=e.data.data.environment_audio,console.log("获取用户信息成功:",x.value)}catch(e){console.error("获取用户信息失败:",e)}},ge=async()=>{var J,Q;const e=p.value.map(r=>{const _e=r.type==="ai"?"面试官":"候选人",we=r.type==="base64"?"[语音回答]":r.text;return`[${r.time}] ${_e}: ${we}`}).join(`
`),a=p.value.filter(r=>r.type==="ai"&&r.text!=="你好，请你先介绍一下自己").length,n=`你是一位专业的${((Q=(J=x.value)==null?void 0:J.data)==null?void 0:Q.role)||"目标岗位"}面试官，正在通过视频+语音形式对一位候选人进行结构化面试。以下是你们的对话记录（包含时间、角色与内容），其中[语音回答]代表候选人的语音回复，请你结合上下文进行判断和提问：

${e}

### 面试官行为准则
1. **专业形象**：语气正式稳重，不使用口语词（如“啦”、“嗯哼”），不使用网络用语
2. **提问风格**：
   - 简洁、清晰、直接
   - 避免解释提问目的
   - 针对候选人已答内容精准追问，如"你刚提到……能展开一下吗？"
   - 多使用开放式问题：如何、为什么、请具体说明等
3. **流程控制**：
   - 当前为第 ${a+1} 个问题，整场控制在 4～7 个问题内
   - 请使用“好的”、“接下来”等自然衔接语过渡
   - 当第7个问题提问完毕后，以“我的问题结束了，你有什么想问我的吗？”作为收尾
4. **深度引导**：
   - 回答模糊时，追问“能举个具体例子吗？”
   - 技术话题可追问原理、实现方案、权衡决策
   - 遇到[语音回答]，请尽量基于上下文提问，但不要直接提到"语音"或"音频"

请严格遵循以上规则，生成下一轮符合专业面试官身份的提问内容，仅返回问题句子本身，无需解释或额外注释。`,l={username:i.value,conversation_id:k.value,audio:S.value.data,text:n};try{const r=await c.post("/api/message/new",l,{headers:{Authorization:s.value?`Bearer ${s.value}`:""}});console.log("AI回复:",r.data),p.value.push({type:"ai",text:r.data.data.text,time:new Date().toLocaleTimeString()}),B.value=!0,f.value=!0,console.log("开始回答"),T.value=r.data.data.text,O(T.value,()=>{f.value=!1}),setTimeout(()=>{g.value&&(g.value.scrollTop=g.value.scrollHeight)},100)}catch(r){console.error("获取AI回复失败:",r);return}},me=async()=>{const e={username:i.value,conversation_id:k.value,text:`你是一个专业的面试分析助手，需要对面试者的视频进行面部表情解析和情感分析。请按以下要求输出内容：

1. **表情解析**：分析面试者的主要面部表情（如微笑、紧张、自信等）及其出现频率。
2. **情感分析**：判断整体情感倾向（积极/中性/消极），并说明依据。
3. **面试表现评判**：基于表情和情感，评估面试者的表现（如是否自然、是否紧张影响表达等）。

**输出要求**：
- 仅返回纯文本，不要使用Markdown符号（如**、-等）。
- 禁止添加额外解释或客套话。

示例输出：
表情解析：面试者微笑频率较高（约60%时间），偶尔出现皱眉（约15%时间）。
情感分析：整体情感倾向积极，因大部分时间保持微笑和眼神接触。
面试表现评判：表现自然，但皱眉可能反映对部分问题感到不确定，建议后续提问时观察其反应。`,video:E.value};try{const a=await c.post("/api/message/new",e,{headers:{Authorization:s.value?`Bearer ${s.value}`:""}});V.value=a.data.data.text,console.log("视频情感分析:",V.value)}catch(a){console.log("视频情感解析失败:",a)}},he=async()=>{const e={username:i.value,conversation_id:k.value,text:`你是一个专业的面试评估助手，用户正在面试的职位是${x.value.data.role},你需对整个面试过程（候选人）进行全面评估，并量化人岗匹配度。请按以下要求输出内容：

1. **候选人表现**：
   - 语言表达（清晰度/逻辑性）
   - 技术能力（回答准确性/深度）
   - 非语言表现（表情/语气）
2. **面试官表现**：
   - 提问质量（相关性/逻辑性）
   - 流程控制（节奏/引导）
3. **人岗匹配度**：
   - 匹配度百分比（如75%）及依据（如技能匹配但经验不足）
4. **改进建议**：
   - 对候选人或面试流程的具体建议

**输出要求**：
- 纯文本，禁止使用Markdown符号（**、-等）
- 严格按以下格式输出，禁止额外内容：

示例输出：
候选人表现：语言表达清晰，技术问题回答准确率85%，但3个深层次问题未答出。面试中保持微笑，但眼神回避频繁。
面试官表现：提问覆盖岗位核心技能，但追问不够深入。面试超时10分钟。
人岗匹配度：70%（核心技能匹配，但项目经验与岗位要求有20%差距）。
改进建议：候选人需补充云计算实战经验；面试官应提前设定问题时间限制。`};try{const a=await c.post("/api/message/new",e,{headers:{Authorization:s.value?`Bearer ${s.value}`:""}});R.value=a.data.data.text,console.log("AI评估回复:",R.value),localStorage.setItem("result",R.value)}catch(a){console.error("获取评估失败:",a)}};return fe(()=>{setTimeout(()=>{j.value=!0},100);const e=localStorage.getItem("token");console.log("进入interviewPage:",e),console.log("当前录音情况",_.value),pe(),H()}),be(()=>{K(),G(),W()}),(e,a)=>(d(),v("div",Ce,[o("div",Re,[o("img",{src:Ie,alt:"",class:M(["stamp",{enter:j.value}])},null,2),a[0]||(a[0]=o("h1",null,"面试模拟",-1)),o("div",Se,[o("span",null,"面试者: "+I(i.value),1)]),o("div",ze,[o("button",{class:"start-button",onClick:oe},I(B.value?"面试已开始":"开始面试"),1),o("button",{class:"back-button",onClick:te},"结束面试")])]),o("div",Ee,[o("div",Ue,[o("div",Ve,[o("div",qe,[o("video",{ref_key:"mainVideo",ref:h,autoplay:"",playsinline:"",class:"mirror"},null,512),a[1]||(a[1]=o("div",{class:"video-placeholder"},null,-1))]),o("div",Le,[o("div",{class:"recording-trigger",onClick:de},[_.value?(d(),v("img",Pe)):(d(),v("img",Me))])])]),o("div",Fe,[a[2]||(a[2]=o("div",{class:"secondary"},[o("img",{src:Be,alt:""})],-1)),o("div",je,[o("p",null,I(U.value||T.value),1)])])]),o("div",Ne,[o("div",Oe,[a[3]||(a[3]=o("h3",null,"历史对话",-1)),o("div",{class:"messages",ref_key:"messagesContainer",ref:g},[(d(!0),v(Ae,null,De(p.value,(n,l)=>(d(),v("div",{key:l,class:M(["message",n.type])},[o("div",We,[n.type==="base64"?(d(),v("div",He,[o("audio",{controls:"",src:n.data},null,8,Ke)])):(d(),v("div",Ge,I(n.text),1)),o("div",Je,I(n.time),1)])],2))),128))],512)]),o("div",Qe,[o("div",{class:M(["ai-indicator",{thinking:f.value}])},null,2),o("span",null,I(f.value?"面试官正在回复...":"面试官已就绪"),1)])]),q.value?(d(),v("div",Xe,a[4]||(a[4]=[o("div",{class:"loading-spinner"},null,-1),o("div",{class:"loading-text"},"正在生成面试报告，请稍候...",-1)]))):ke("",!0)])]))}},aa=ye(Ye,[["__scopeId","data-v-ae66cb6b"]]);export{aa as default};
