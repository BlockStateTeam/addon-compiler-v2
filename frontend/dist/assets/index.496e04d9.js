(function(){const t=document.createElement("link").relList;if(t&&t.supports&&t.supports("modulepreload"))return;for(const o of document.querySelectorAll('link[rel="modulepreload"]'))i(o);new MutationObserver(o=>{for(const n of o)if(n.type==="childList")for(const r of n.addedNodes)r.tagName==="LINK"&&r.rel==="modulepreload"&&i(r)}).observe(document,{childList:!0,subtree:!0});function a(o){const n={};return o.integrity&&(n.integrity=o.integrity),o.referrerpolicy&&(n.referrerPolicy=o.referrerpolicy),o.crossorigin==="use-credentials"?n.credentials="include":o.crossorigin==="anonymous"?n.credentials="omit":n.credentials="same-origin",n}function i(o){if(o.ep)return;o.ep=!0;const n=a(o);fetch(o.href,n)}})();function v(e){return window.go.main.App.CompilePack(e)}function I(e){return window.go.main.App.GetData(e)}function S(e){return window.go.main.App.GetImage(e)}function x(){return window.go.main.App.GetUserSetting()}function k(e,t){return window.go.main.App.Notify(e,t)}function y(e){return window.go.main.App.NotifyText(e)}function P(e){return window.go.main.App.OpenDirectoryDialog(e)}function C(e){return window.go.main.App.SaveUserSetting(e)}function L(){return window.go.main.App.SelfVersion()}function M(e,t,a){return window.go.main.App.UpdateScriptVersion(e,t,a)}function T(){window.runtime.WindowReload()}function N(e){window.runtime.BrowserOpenURL(e)}const w={default:{"--background-color":"#242424","--text-color":"#fff","--font-size":"1.5vw","--table-width":"50vw","--table-background":"#171717","--button-background":"#666","--button-hover-background":"#393939","--button-active-background":"#000000","--header-font-size":"4vw","--h2-font-size":"2.5vw","--h3-font-size":"2.75vw","--h3-hover-font-size":"3vw","--h3-hover-color":"#80a2ff","--h4-font-size":"1.5vw","--select-background":"#393939","--select-text-color":"#fff","--default-font":"Arial, sans-serif","--h1-font":"Arial, sans-serif"},light:{"--background-color":"#ffffff","--text-color":"#242424","--font-size":"1.5vw","--table-width":"50vw","--table-background":"#d5d5d5","--button-background":"#cacaca","--button-hover-background":"#ffffff","--button-active-background":"#c6c6c6","--header-font-size":"4vw","--h2-font-size":"2.5vw","--h3-font-size":"2.75vw","--h3-hover-font-size":"3vw","--h3-hover-color":"#1654ff","--h4-font-size":"1.5vw","--select-background":"#393939","--select-text-color":"#fff","--default-font":"Arial, sans-serif","--h1-font":"Arial, sans-serif"},minecraft:{"--background-color":"#242424","--text-color":"#fff","--font-size":"1.3vw","--table-width":"50vw","--table-background":"#171717","--button-background":"#666","--button-hover-background":"#393939","--button-active-background":"#c9c9c9","--header-font-size":"4vw","--h2-font-size":"2.5vw","--h3-font-size":"2.75vw","--h3-hover-font-size":"3vw","--h3-hover-color":"#80a2ff","--h4-font-size":"1.3vw","--select-background":"#393939","--select-text-color":"#fff","--default-font":"Minecraft","--h1-font":"MinecraftTen"}};let d,g=!0,$=["BUTTON","SELECT","SPAN","TD"],l=document.getElementById("contextMenu");function O(e){return e=e.charAt(0).toUpperCase()+e.slice(1),e.replace(/([a-z])([A-Z])/g,"$1 $2")}document.addEventListener("mouseover",e=>{let t=e.target;if($.includes(t.tagName))switch(t.tagName){case"TD":let a=t.closest("tr"),i=a.getAttribute("data-current-version")==="null-script"?"":" with BetaAPI version "+a.getAttribute("data-current-version");l.innerHTML=`<strong>${O(a.getAttribute("data-pack-type"))}</strong> ${i}`;break;case"BUTTON":let o=t.textContent;if(o==="Update API")l.innerHTML="Click to <strong>update</strong> the Script API version for this pack";else if(o==="Compile"){let n=t.closest("tr"),r=document.getElementById("saveMode").value==="desktop"?"to Desktop":"";l.innerHTML=`Click to compile <strong>${n.getAttribute("data-pack-name")}</strong> as ${document.getElementById("saveFormat").value.toUpperCase()} format ${r}`}else o==="Dev Mode"&&(l.innerHTML="Click to toggle between <strong>Development</strong> and <strong>Release</strong> mode");break;case"SELECT":t.id==="saveFormat"?l.innerHTML="Change export <strong>format</strong> for compiled packs":t.id==="saveMode"?l.innerHTML="Change export <strong>location</strong> for compiled packs":l.innerHTML="Change <strong>theme</strong> for the app";break;case"SPAN":t.id==="latestScriptVersion"&&(l.innerHTML="Script API version refer to <strong>BetaAPI</strong> experimental features.<br><br> Version: "+t.textContent);break}else l.innerHTML="Hover above buttons or selection <br> to get more info <br><br> Click to hide this pannel"});document.addEventListener("mousemove",e=>{let t=window.screen.width-(e.clientX+425),a=window.screen.height-(e.clientY+document.getElementById("contextMenu").offsetHeight+75);t<0?t=e.clientX-425:t=e.clientX+5,a<0?a=e.clientY-document.getElementById("contextMenu").offsetHeight-5:a=e.clientY+5,l.style.top=`${a}px`,l.style.left=`${t}px`});let f=!1;document.addEventListener("contextmenu",e=>{if(e.preventDefault(),f){document.getElementById("contextMenu").style.opacity="0",f=!1;return}f=!0;let t=window.screen.width-(e.clientX+425);t<0?t=e.clientX-420:t=e.clientX,l.style.display="block",l.style.opacity="1",document.addEventListener("click",()=>{f=!1,document.getElementById("contextMenu").style.opacity="0"})});function E(e){const t=document.documentElement;for(const[a,i]of Object.entries(e))t.style.setProperty(a,i)}L().then(e=>{document.getElementById("selfVersion").textContent="v"+e});x().then(e=>{d=JSON.parse(e),document.getElementById("saveFormat").value=d.format||"mcaddon",document.getElementById("saveMode").value=d.exportLocation||"desktop",document.getElementById("uiselect").value=d.theme||"default",document.getElementById("modeButton").getAttribute("data-mode")===d.mode||h(),E(w[d.theme])});function A(e){let t="";e==="development"&&(t="Development "),document.getElementById("resourcePack")&&(document.getElementById("resourcePack").textContent=t+"Resource Packs",document.getElementById("addOnPack").textContent=t+"Add-On Packs",document.getElementById("behaviorPack").textContent=t+"Behavior Packs")}function m(e,t){d[e]=t,C(JSON.stringify(d))}function h(){if(g){console.log("Toggling Mode");let e=document.getElementById("modeButton");e.getAttribute("data-mode")==="release"?(e.setAttribute("data-mode","development"),e.textContent="Dev Mode"):(e.setAttribute("data-mode","release"),e.textContent="Release Mode");try{m("mode",e.getAttribute("data-mode")),b()}catch(t){console.error(t)}A(e.getAttribute("data-mode"))}else setTimeout(h,100)}document.getElementById("modeButton").addEventListener("click",e=>{h()});document.getElementById("saveFormat").addEventListener("change",e=>{m("format",e.target.value)});document.getElementById("saveMode").addEventListener("change",e=>{m("exportLocation",e.target.value)});document.getElementById("uiselect").addEventListener("change",e=>{let t=e.target.value;m("theme",t),E(w[t])});async function z(e){let t=["resourcePack","addOnPack","behaviorPack"];const a=[];for(let i=0;i<e.length;i++){let o=e[i],n=`<h2 id="${t[i]}">Error</h2><table><tbody>`;if(o)for(let r=0;r<o.length;r++){let c=o[r],s="null-script";c.ScriptState&&(s=c.ScriptState);let u=await S(c.Path||c.ResourcePath);n+=`
                    <tr data-current-version="${s}" data-pack-name="${c.CleanName}" data-pack-type="${t[i]}" data-pack-path="${c.Path}" data-pack-rp-path="${c.ResourcePath}" data-pack-bp-path="${c.BehaviorPath}">
                        <td><img class="projectImage" src="${u}"></td>
                        <td id="textCell">${c.CleanName}</td>
                        <td id="buttonCell1" class="updateScriptVersionButton"><button class="updateVersion" style="display:none;">Update API</button></td>
                        <td id="buttonCell2"><button class="compileButton">Compile</button></td>
                    </tr>`}n+="</tbody></table>",a.push(n)}return a.join("")}function D(){let e=[];fetch("https://registry.npmjs.org/@minecraft/server").then(t=>{if(!t.ok)throw y("Failed to check for Script API version. Please check your internet connection."),new Error("Network response was not ok");return t.json()}).then(t=>{try{const a=Math.max(...Object.keys(t.time).filter(n=>/^\d+\.\d+\.\d+/.test(n)).filter(n=>n.endsWith("-stable")).map(n=>{const r=n.match(/^(\d+\.\d+\.\d+)/);return r?r[1].split(".").reduce((c,s,u)=>c+s*Math.pow(100,2-u),0):null})).toString();let i=a.length;for(;i>0;i-=2)e.unshift(a.slice(Math.max(i-2,0),i));let o=e.map(n=>Number(n)).join(".")+"-beta";document.getElementById("latestScriptVersion").innerText=o,document.getElementById("latestScriptContainer").style.visibility="visible",[...document.getElementsByClassName("updateScriptVersionButton")].forEach(n=>{let r=n.closest("tr").getAttribute("data-current-version");r!=="null-script"&&r!==o&&(n.querySelector("button").style.display="inline-block")})}catch(a){console.error("Error: "+a)}}).catch(t=>{console.warn(t)})}function p(e){return e==="undefined"?void 0:e}function b(){g=!1;let e=document.getElementById("modeButton").getAttribute("data-mode")==="development";console.log(e),I(e).then(t=>{const a=JSON.parse(t);console.log(a),z(a).then(i=>{g=!0,document.getElementById("table").innerHTML=i,A(document.getElementById("modeButton").getAttribute("data-mode")),D(),document.querySelectorAll("button").forEach(o=>{o.addEventListener("click",()=>{let n=o.closest("tr");if(o.classList.contains("updateVersion")){let r=document.getElementById("latestScriptVersion").textContent,c=p(n.getAttribute("data-pack-bp-path"))||n.getAttribute("data-pack-path");M(c,n.getAttribute("data-current-version"),r),b(),y(`Updating script for ${n.getAttribute("data-pack-name")}`)}else if(o.classList.contains("compileButton")){let r=n.getAttribute("data-pack-name"),c=n.getAttribute("data-pack-path")||n.getAttribute("data-pack-rp-path");c+="\\pack_icon.png";let s={CleanName:n.getAttribute("data-pack-name"),PackType:n.getAttribute("data-pack-type"),ResoucePackPath:p(n.getAttribute("data-pack-rp-path"))||n.getAttribute("data-pack-path"),BehaviorPackPath:p(n.getAttribute("data-pack-bp-path"))||n.getAttribute("data-pack-path"),ExportPath:"desktop",Format:document.getElementById("saveFormat").value};console.log(s);let u=document.getElementById("saveFormat").value;document.getElementById("saveMode").value==="choose"?P({defaultDirectory:"C:/",title:"Select Export Directory"}).then(B=>{s.ExportPath=B,v(JSON.stringify(s)).then(()=>{k(`Finished compiling: 
 ${r}.${u}`,c)})}).catch():v(JSON.stringify(s)).then(()=>{k(`Finished compiling: 
 ${r}.${u}`,c)})}})}),document.getElementById("expand").addEventListener("click",()=>{N("https://blockstate.team")})})})}b();window.addEventListener("online",()=>{console.log("Back Online!"),T()});
