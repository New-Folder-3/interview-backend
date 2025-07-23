import{_ as fe,r as t,o as Ie,a as be,c as p,b as o,n as j,d as ke,t as k,e as Be,f as Ae,F as De,g as xe,h as u,u as Te,i as g}from"./index-BfBhGJYP.js";const Re="/assets/%E5%A3%B0%E9%9F%B3%E5%85%B3%E9%97%AD-DpSzKatq.png",Se="/assets/%E5%BC%80%E5%90%AF%E5%A3%B0%E9%9F%B3-qtu0_L9e.png",$e={class:"interview-container"},Ce={class:"header"},ze={class:"info"},Ee={class:"header-buttons"},Ue={class:"main-content"},Ve={class:"left"},qe={class:"video-panel"},Le={class:"video-main"},Me={class:"recording-panel"},Pe={key:0,src:Re,alt:"开始录音",class:"record-img"},Fe={key:1,src:Se,alt:"停止录音",class:"record-img"},je={class:"side-panel"},Ne={class:"ai-message-box"},Oe={class:"interaction-panel"},We={class:"conversation-area"},He={class:"message-content"},Ke={key:0,class:"audio-player"},Ge=["src"],Je={key:1,class:"message-text"},Qe={class:"message-time"},Xe={class:"ai-status"},Ye={key:0,class:"loading-overlay"},Ze={__name:"InterviewPage",setup(ea){const Z=Te(),R=t(null),B=t(null),S=t(null),s=t(localStorage.getItem("token")||""),i=t(localStorage.getItem("username")||""),c=t(!1),A=t(!1),v=t([]),m=t(null),ee=t(!1),z=t(""),d=t(""),D=t([]),h=t(null),_=t(null),w=t(!1),x=t(null),E=t(null),U=t(null),y=t(null),I=t([]),N=t(!1),$=t(null),V=t(null),b=t(""),q=t(""),O=t(null),C=t(""),L=t(""),M=t(!1),W=t(!1),P=t(!1);t(null);const f=t(null),H=t(!1),ae=async()=>{console.log("当前用户名:",i.value),console.log("当前token:",s.value);try{const e={username:i.value,prefer_role:0};console.log("发送的请求数据:",e);const a=await u.post("/api/conversation/new",e,{headers:{Authorization:s?`Bearer ${s.value}`:""}});B.value=a.data.data.conversation_id,console.log("成功创建对话，ID:",B.value)}catch(e){console.error("创建会话失败:",e)}},te=async()=>{const e={username:i.value};try{const a=await u.post("api/interview/new",e,{headers:{Authorization:s.value?`Bearer ${s.value}`:""}});console.log("获取面试ID成功:",a.data.data),S.value=a.data.data,localStorage.setItem("interview_id",S.value)}catch(a){console.error("获取面试ID失败:",a)}},oe=async()=>{const e={username:i.value,interview_id:S.value,interview_conversation:C.value,video_url:V.value,emotion_conversation:L.value};try{const a=await u.post("api/interview/update",e,{headers:{Authorization:s.value?`Bearer ${s.value}`:""}});console.log("更新面试数据成功:",a.data)}catch(a){console.error("更新面试数据失败:",a)}},se=async()=>{console.log("面试结束，进入结果界面"),M.value=!0;try{J(),K(),await we(),await _e(),await oe(),Z.push(`/result/${S.value}`)}catch(e){console.error("跳转结果页失败:",e),M.value=!1,alert("生成报告失败，请重试")}},F=(e,a)=>{q.value="";let n=0;const r=setInterval(()=>{q.value+=e[n],n++,n>=e.length&&(clearInterval(r),a&&a())},80)},ne=async()=>{if(A.value){alert("面试已开始"),z.value="";return}await G(),console.log("是否开启白噪音：",P.value),P.value&&le(),te(),ue(),A.value=!0,c.value=!0,console.log("开始面试"),await ae();const e="你好，请你先介绍一下自己";b.value="",F(e,()=>{c.value=!1,b.value=e}),setTimeout(()=>{m.value&&(m.value.scrollTop=m.value.scrollHeight,v.value.push({type:"ai",text:e,time:new Date().toLocaleTimeString()}))},1e3)},le=()=>{f.value||(f.value=new Audio("../assets/audio/noisy.mp3"),f.value.loop=!0),f.value.play().then(()=>{H.value=!0,console.log("白噪音开始播放")}).catch(e=>{console.error("白噪音播放失败:",e)})},K=()=>{f.value&&(f.value.pause(),f.value.currentTime=0,H.value=!1,console.log("白噪音已停止"))},G=async()=>{try{const e=await navigator.mediaDevices.getUserMedia({video:!0,audio:!0});U.value=e,_.value&&(_.value.srcObject=e)}catch(e){console.error("无法访问摄像头:",e),alert("无法访问摄像头，请检查权限。")}},re=()=>{U.value&&(I.value=[],y.value=new MediaRecorder(U.value,{mimeType:"video/webm; codecs=vp9"}),y.value.ondataavailable=e=>{e.data.size>0&&I.value.push(e.data)},y.value.onstop=()=>{const e=new Blob(I.value,{type:"video/webm"});$.value=URL.createObjectURL(e),console.log("Video recording stopped",$.value)},y.value.start(),N.value=!0,console.log("Video recording started"))},ie=()=>{y.value&&y.value.state==="recording"&&(y.value.stop(),N.value=!1,console.log("Video recording stopped"))},ue=()=>{console.log("开始录制30s视频用于情感分析"),re(),O.value=setTimeout(async()=>{await ie(),await new Promise(e=>{const a=setInterval(()=>{$.value&&(clearInterval(a),e())},100)}),ce(),O.value=null},1e4)},ce=async()=>{if(console.log("开始保存视频片段"),!$.value){console.error("视频URL不可用");return}try{if(!I.value||I.value.length===0){console.error("没有可用的视频片段");return}const e=new FormData,a=new Blob(I.value,{type:"video/webm"});console.log("videotype",a.type);const n={file:a};e.append("video",a,"30-second-clip.webm");const r=await u.post("/api/upload/video",n,{headers:{Authorization:s.value?`Bearer ${s.value}`:"","Content-Type":"multipart/form-data"}});console.log("视频上传成功:",r.data),V.value=r.data.data}catch(e){console.error("视频上传失败:",e)}},J=()=>{_.value&&_.value.srcObject&&(_.value.srcObject.getTracks().forEach(e=>e.stop()),_.value.srcObject=null)},ve=e=>new Promise((a,n)=>{const r=new FileReader;r.onloadend=()=>{a(r.result)},r.onerror=n,r.readAsDataURL(e)}),de=async()=>{if(!d.value){console.error("没有音频数据可上传");return}try{const e=await u.post("api/upload/audio64",{base64:d.value},{headers:{Authorization:s.value?`Bearer ${s.value}`:""}});E.value=e.data,console.log("音频上传成功:",E.value.data)}catch(e){console.error("Error uploading audio:",e)}},pe=async()=>{if(A.value)try{x.value=await navigator.mediaDevices.getUserMedia({audio:!0}),D.value=[],h.value=new MediaRecorder(x.value,{mimeType:"audio/webm"}),h.value.ondataavailable=e=>{e.data.size>0&&D.value.push(e.data)},h.value.onstop=async()=>{const e=new Blob(D.value,{type:"audio/webm"});d.value=await ve(e),console.log("音频转换为Base64:",d.value),D.value=[]},h.value.start(),w.value=!0}catch(e){w.value=!1,console.error("无法访问麦克风或不支持录音:",e)}else alert("请先开始面试")},Q=()=>{h.value&&h.value.state!=="inactive"&&h.value.stop(),x.value&&(x.value.getTracks().forEach(e=>e.stop()),x.value=null),w.value=!1},ge=async()=>{if(w.value){if(Q(),console.log("停止录音"),await new Promise(e=>{const a=setInterval(()=>{!w.value&&d.value&&(clearInterval(a),e())},100)}),d.value){v.value=[...v.value,{type:"base64",data:d.value,time:new Date().toLocaleTimeString()}],ee.value=!0;try{await de(),await he()}catch(e){console.error("处理音频时出错:",e)}d.value=""}}else console.log("开始录音"),D.value=[],await pe()},me=async()=>{try{const e=await u.get("api/user/get",{headers:{Authorization:`Bearer ${s.value}`},params:{username:i.value}});R.value=e.data,P.value=e.data.data.environment_audio,console.log("获取用户信息成功:",R.value)}catch(e){console.error("获取用户信息失败:",e)}},he=async()=>{var X,Y;const e=v.value.map(l=>{const T=l.type==="ai"?"面试官":"候选人",ye=l.type==="base64"?"[语音回答]":l.text;return`[${l.time}] ${T}: ${ye}`}).join(`
`),a=v.value.filter(l=>l.type==="ai"&&l.text!=="你好，请你先介绍一下自己").length,n=`你是一位专业的${((Y=(X=R.value)==null?void 0:X.data)==null?void 0:Y.role)||"目标岗位"}面试官，正在通过视频+语音形式对一位候选人进行结构化面试。以下是你们的对话记录（包含时间、角色与内容），其中[语音回答]代表候选人的语音回复，请你结合上下文进行判断和提问：

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

请严格遵循以上规则，生成下一轮符合专业面试官身份的提问内容，仅返回问题句子本身，无需解释或额外注释。`,r={username:i.value,conversation_id:B.value,audio:E.value.data,text:n};try{const l=await u.post("/api/message/new",r,{headers:{Authorization:s.value?`Bearer ${s.value}`:""}});if(l.data.data.text===z.value){const T="本次面试到此结束，请结束面试查看分析结果";v.value.push({type:"ai",text:T,time:new Date().toLocaleTimeString()}),b.value=T,c.value=!0,F(T,()=>{c.value=!1});return}z.value=l.data.data.text,v.value.push({type:"ai",text:l.data.data.text,time:new Date().toLocaleTimeString()}),A.value=!0,c.value=!0,b.value=l.data.data.text,F(b.value,()=>{c.value=!1}),setTimeout(()=>{m.value&&(m.value.scrollTop=m.value.scrollHeight)},100)}catch(l){console.error("获取AI回复失败:",l);return}},_e=async()=>{const e={username:i.value,conversation_id:B.value,text:`你是一个专业的面试分析助手，需要对面试者的视频进行面部表情解析和情感分析。请按以下要求输出内容：

1. **表情解析**：分析面试者的主要面部表情（如微笑、紧张、自信等）及其出现频率。
2. **情感分析**：判断整体情感倾向（积极/中性/消极），并说明依据。
3. **面试表现评判**：基于表情和情感，评估面试者的表现（如是否自然、是否紧张影响表达等）。

**输出要求**：
- 仅返回纯文本，不要使用Markdown符号（如**、-等）。
- 禁止添加额外解释或客套话。

示例输出：
表情解析：面试者微笑频率较高（约60%时间），偶尔出现皱眉（约15%时间）。
情感分析：整体情感倾向积极，因大部分时间保持微笑和眼神接触。
面试表现评判：表现自然，但皱眉可能反映对部分问题感到不确定，建议后续提问时观察其反应。`,video:V.value};try{const a=await u.post("/api/message/new",e,{headers:{Authorization:s.value?`Bearer ${s.value}`:""}});L.value=a.data.data.text,console.log("视频情感分析:",L.value)}catch(a){console.log("视频情感解析失败:",a)}},we=async()=>{const e={username:i.value,conversation_id:B.value,text:`你是一个专业的面试评估助手，用户正在面试的职位是${R.value.data.role},你需对整个面试过程（候选人）进行全面评估，并量化人岗匹配度。请按以下要求输出内容：

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
改进建议：候选人需补充云计算实战经验；面试官应提前设定问题时间限制。`};try{const a=await u.post("/api/message/new",e,{headers:{Authorization:s.value?`Bearer ${s.value}`:""}});C.value=a.data.data.text,console.log("AI评估回复:",C.value),localStorage.setItem("result",C.value)}catch(a){console.error("获取评估失败:",a)}};return Ie(()=>{setTimeout(()=>{W.value=!0},100);const e=localStorage.getItem("token");console.log("进入interviewPage:",e),console.log("当前录音情况",w.value),me(),G()}),be(()=>{J(),Q(),K()}),(e,a)=>(g(),p("div",$e,[o("div",Ce,[o("img",{src:ke,alt:"",class:j(["stamp",{enter:W.value}])},null,2),a[0]||(a[0]=o("h1",null,"面试模拟",-1)),o("div",ze,[o("span",null,"面试者: "+k(i.value),1)]),o("div",Ee,[o("button",{class:"start-button",onClick:ne},k(A.value?"面试已开始":"开始面试"),1),o("button",{class:"back-button",onClick:se},"结束面试")])]),o("div",Ue,[o("div",Ve,[o("div",qe,[o("div",Le,[o("video",{ref_key:"mainVideo",ref:_,autoplay:"",playsinline:"",class:"mirror"},null,512),a[1]||(a[1]=o("div",{class:"video-placeholder"},null,-1))]),o("div",Me,[o("div",{class:"recording-trigger",onClick:ge},[w.value?(g(),p("img",Fe)):(g(),p("img",Pe))])])]),o("div",je,[a[2]||(a[2]=o("div",{class:"secondary"},[o("img",{src:Ae,alt:""})],-1)),o("div",Ne,[o("p",null,k(q.value||b.value),1)])])]),o("div",Oe,[o("div",We,[a[3]||(a[3]=o("h3",null,"历史对话",-1)),o("div",{class:"messages",ref_key:"messagesContainer",ref:m},[(g(!0),p(De,null,xe(v.value,(n,r)=>(g(),p("div",{key:r,class:j(["message",n.type])},[o("div",He,[n.type==="base64"?(g(),p("div",Ke,[o("audio",{controls:"",src:n.data},null,8,Ge)])):(g(),p("div",Je,k(n.text),1)),o("div",Qe,k(n.time),1)])],2))),128))],512)]),o("div",Xe,[o("div",{class:j(["ai-indicator",{thinking:c.value}])},null,2),o("span",null,k(c.value?"面试官正在回复...":"面试官已就绪"),1)])]),M.value?(g(),p("div",Ye,a[4]||(a[4]=[o("div",{class:"loading-spinner"},null,-1),o("div",{class:"loading-text"},"正在生成面试报告，请稍候...",-1)]))):Be("",!0)])]))}},ta=fe(Ze,[["__scopeId","data-v-a3498fc7"]]);export{ta as default};
