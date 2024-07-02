(function(){const t=document.createElement("link").relList;if(t&&t.supports&&t.supports("modulepreload"))return;for(const o of document.querySelectorAll('link[rel="modulepreload"]'))l(o);new MutationObserver(o=>{for(const n of o)if(n.type==="childList")for(const r of n.addedNodes)r.tagName==="LINK"&&r.rel==="modulepreload"&&l(r)}).observe(document,{childList:!0,subtree:!0});function a(o){const n={};return o.integrity&&(n.integrity=o.integrity),o.referrerpolicy&&(n.referrerPolicy=o.referrerpolicy),o.crossorigin==="use-credentials"?n.credentials="include":o.crossorigin==="anonymous"?n.credentials="omit":n.credentials="same-origin",n}function l(o){if(o.ep)return;o.ep=!0;const n=a(o);fetch(o.href,n)}})();function I(e){return window.go.main.App.AutoUpdate(e)}function k(e){return window.go.main.App.CompilePack(e)}function S(e){return window.go.main.App.GetData(e)}function P(e){return window.go.main.App.GetImage(e)}function C(){return window.go.main.App.GetUserSetting()}function T(e,t){return window.go.main.App.Normalize(e,t)}function y(e,t){return window.go.main.App.Notify(e,t)}function p(e){return window.go.main.App.NotifyText(e)}function z(e){return window.go.main.App.OpenDirectoryDialog(e)}function L(e){return window.go.main.App.SaveUserSetting(e)}function M(){return window.go.main.App.SelfVersion()}function N(e,t,a){return window.go.main.App.UpdateScriptVersion(e,t,a)}function O(e,t,a){return window.runtime.EventsOnMultiple(e,t,a)}function w(e,t){return O(e,t,-1)}function $(){window.runtime.WindowReload()}function V(e){window.runtime.BrowserOpenURL(e)}const E={default:{"--background-normalize-panel":"rgba(0, 0, 0, 0.5)","--background-color":"#242424","--grid-color":"#ffffff00","--text-color":"#fff","--font-size":"1.5vw","--table-width":"60vw","--table-background":"#171717","--button-background":"#666","--button-hover-background":"#393939","--button-active-background":"#000000","--header-font-size":"4vw","--h2-font-size":"2.5vw","--h3-font-size":"2.75vw","--h3-hover-font-size":"3vw","--h3-hover-color":"#80a2ff","--h4-font-size":"1.5vw","--select-background":"#393939","--select-text-color":"#fff","--default-font":"Arial, sans-serif","--h1-font":"Arial, sans-serif"},light:{"--background-normalize-panel":"rgba(255, 255, 255, 0.5)","--background-color":"#ffffff","--grid-color":"#00000000","--text-color":"#242424","--font-size":"1.5vw","--table-width":"60vw","--table-background":"#d5d5d5","--button-background":"#cacaca","--button-hover-background":"#ffffff","--button-active-background":"#c6c6c6","--header-font-size":"4vw","--h2-font-size":"2.5vw","--h3-font-size":"2.75vw","--h3-hover-font-size":"3vw","--h3-hover-color":"#1654ff","--h4-font-size":"1.5vw","--select-background":"#393939","--select-text-color":"#fff","--default-font":"Arial, sans-serif","--h1-font":"Arial, sans-serif"},minecraft:{"--background-normalize-panel":"rgba(0, 0, 0, 0.5)","--background-color":"#242424","--grid-color":"#ffffff00","--text-color":"#fff","--font-size":"1.3vw","--table-width":"60vw","--table-background":"#171717","--button-background":"#666","--button-hover-background":"#393939","--button-active-background":"#c9c9c9","--header-font-size":"4vw","--h2-font-size":"2.5vw","--h3-font-size":"2.75vw","--h3-hover-font-size":"3vw","--h3-hover-color":"#80a2ff","--h4-font-size":"1.3vw","--select-background":"#393939","--select-text-color":"#fff","--default-font":"Minecraft","--h1-font":"MinecraftTen"},grid:{"--background-normalize-panel":"rgba(0, 0, 0, 0.5)","--background-color":"#242424","--grid-color":"#ffffff10","--text-color":"#fff","--font-size":"1.5vw","--table-width":"60vw","--table-background":"#171717","--button-background":"#666","--button-hover-background":"#393939","--button-active-background":"#000000","--header-font-size":"4vw","--h2-font-size":"2.5vw","--h3-font-size":"2.75vw","--h3-hover-font-size":"3vw","--h3-hover-color":"#80a2ff","--h4-font-size":"1.5vw","--select-background":"#393939","--select-text-color":"#fff","--default-font":"Arial, sans-serif","--h1-font":"Arial, sans-serif"}};let d,h=!0,D=["BUTTON","SELECT","SPAN","TD"],c=document.getElementById("contextMenu");function U(e){return e=e.charAt(0).toUpperCase()+e.slice(1),e.replace(/([a-z])([A-Z])/g,"$1 $2")}document.addEventListener("mouseover",e=>{let t=e.target;if(D.includes(t.tagName))switch(t.tagName){case"TD":let a=t.closest("tr"),l=a.getAttribute("data-current-version")==="null-script"?"":" with BetaAPI version "+a.getAttribute("data-current-version");c.innerHTML=`<strong>${U(a.getAttribute("data-pack-type"))}</strong> ${l}`;break;case"BUTTON":let o=t.textContent;if(o==="Update API")c.innerHTML="Click to <strong>update</strong> the Script API version for this pack";else if(o==="Compile"){let n=t.closest("tr"),r=document.getElementById("saveMode").value==="desktop"?"to Desktop":"";c.innerHTML=`Click to compile <strong>${n.getAttribute("data-pack-name")}</strong> as ${document.getElementById("saveFormat").value.toUpperCase()} format ${r}`}else o==="Dev Mode"&&(c.innerHTML="Click to toggle between <strong>Development</strong> and <strong>Release</strong> mode");break;case"SELECT":t.id==="saveFormat"?c.innerHTML="Change export <strong>format</strong> for compiled packs":t.id==="saveMode"?c.innerHTML="Change export <strong>location</strong> for compiled packs":c.innerHTML="Change <strong>theme</strong> for the app";break;case"SPAN":t.id==="latestScriptVersion"&&(c.innerHTML="Script API version refer to <strong>BetaAPI</strong> experimental features.<br><br> Version: "+t.textContent);break}else c.innerHTML="Hover above buttons or selection <br> to get more info <br><br> Click to hide this pannel"});document.addEventListener("mousemove",e=>{let t=window.screen.width-(e.clientX+425),a=window.screen.height-(e.clientY+document.getElementById("contextMenu").offsetHeight+75);t<0?t=e.clientX-425:t=e.clientX+5,a<0?a=e.clientY-document.getElementById("contextMenu").offsetHeight-5:a=e.clientY+5,c.style.top=`${a}px`,c.style.left=`${t}px`});let f=!1;document.addEventListener("contextmenu",e=>{if(e.preventDefault(),f){document.getElementById("contextMenu").style.opacity="0",f=!1;return}f=!0;let t=window.screen.width-(e.clientX+425);t<0?t=e.clientX-420:t=e.clientX,c.style.display="block",c.style.opacity="1",document.addEventListener("click",()=>{f=!1,document.getElementById("contextMenu").style.opacity="0"})});function A(e){const t=document.documentElement;for(const[a,l]of Object.entries(e))t.style.setProperty(a,l)}M().then(e=>{document.getElementById("selfVersion").textContent="v"+e,fetch("https://compiler.blockstate.team/index.json").then(t=>{if(!t.ok)throw new Error("Network response was not ok "+t.statusText);return t.json()}).then(t=>{let a=t.version.split("."),l=e.split(".");+a[0]*1e6+ +a[1]*1e3+ +a[2]>+l[0]*1e6+ +l[1]*1e3+ +l[2]&&(console.log("New Version Available"),I(t.downloadUrl).then(o=>{console.log(o)}))}).catch(t=>{console.error("There was a problem with the fetch operation:",t)})});C().then(e=>{d=JSON.parse(e),document.getElementById("saveFormat").value=d.format||"mcaddon",document.getElementById("saveMode").value=d.exportLocation||"desktop",document.getElementById("uiselect").value=d.theme||"default",document.getElementById("modeButton").getAttribute("data-mode")===d.mode||b(),A(E[d.theme])});function B(e){let t="";e==="development"&&(t="Development "),document.getElementById("resourcePack")&&(document.getElementById("resourcePack").textContent=t+"Resource Packs",document.getElementById("addOnPack").textContent=t+"Add-On Packs",document.getElementById("behaviorPack").textContent=t+"Behavior Packs")}function m(e,t){d[e]=t,L(JSON.stringify(d))}function b(){if(h){console.log("Toggling Mode");let e=document.getElementById("modeButton");e.getAttribute("data-mode")==="release"?(e.setAttribute("data-mode","development"),e.textContent="Dev Mode"):(e.setAttribute("data-mode","release"),e.textContent="Release Mode");try{m("mode",e.getAttribute("data-mode")),v()}catch(t){console.error(t)}B(e.getAttribute("data-mode"))}else setTimeout(b,100)}document.getElementById("modeButton").addEventListener("click",e=>{b()});document.getElementById("saveFormat").addEventListener("change",e=>{m("format",e.target.value)});document.getElementById("saveMode").addEventListener("change",e=>{m("exportLocation",e.target.value)});document.getElementById("uiselect").addEventListener("change",e=>{let t=e.target.value;m("theme",t),A(E[t])});async function H(e){let t=["resourcePack","addOnPack","behaviorPack"];const a=[];for(let l=0;l<e.length;l++){let o=e[l],n=`<h2 id="${t[l]}">Error</h2><table><tbody>`;if(o)for(let r=0;r<o.length;r++){let i=o[r],s="null-script";i.ScriptState&&(s=i.ScriptState);let u=await P(i.Path||i.ResourcePath);n+=`
                    <tr data-current-version="${s}" data-pack-name="${i.CleanName}" data-pack-type="${t[l]}" data-pack-path="${i.Path}" data-pack-rp-path="${i.ResourcePath}" data-pack-bp-path="${i.BehaviorPath}" data-pack-signatures="${i.IsSignatures}">
                        <td class="noneData"><img class="projectImage" src="${u}"></td>
                        <td id="textCell" class="averageTextCell">${i.CleanName}</td>
                        <td id="buttonCell1" class="updateScriptVersionButton"><button class="updateVersion" style="visibility:hidden;">Update API</button></td>
                        <td id="buttonCell2"><button class="compileButton">Compile</button></td>
                        <td id="buttonCell3" class="normalizeButton"><button class="normalizePack" style="visibility:hidden;">Normalize Pack</button></td>
                    </tr>`}n+="</tbody></table>",a.push(n)}return a.join("")}function j(){let e=[];fetch("https://registry.npmjs.org/@minecraft/server").then(t=>{if(!t.ok)throw p("Failed to check for Script API version. Please check your internet connection."),new Error("Network response was not ok");return t.json()}).then(t=>{try{const a=Math.max(...Object.keys(t.time).filter(n=>/^\d+\.\d+\.\d+/.test(n)).filter(n=>n.endsWith("-stable")).map(n=>{const r=n.match(/^(\d+\.\d+\.\d+)/);return r?r[1].split(".").reduce((i,s,u)=>i+s*Math.pow(100,2-u),0):null})).toString();let l=a.length;for(;l>0;l-=2)e.unshift(a.slice(Math.max(l-2,0),l));let o=e.map(n=>Number(n)).join(".")+"-beta";document.getElementById("latestScriptVersion").innerText=o,document.getElementById("latestScriptContainer").style.visibility="visible",[...document.getElementsByClassName("normalizeButton")].forEach(n=>{let r=n.closest("tr").getAttribute("data-pack-signatures")==="true";r&&(console.log(r),n.querySelector("button").style.visibility="visible")}),[...document.getElementsByClassName("updateScriptVersionButton")].forEach(n=>{let r=n.closest("tr").getAttribute("data-current-version");r!=="null-script"&&r.endsWith("-beta")&&r!==o&&(n.querySelector("button").style.visibility="visible")})}catch(a){console.error("Error: "+a)}}).catch(t=>{console.warn(t)})}function g(e){return e==="undefined"?void 0:e}function v(){h=!1;let e=document.getElementById("modeButton").getAttribute("data-mode")==="development";console.log(e),S(e).then(t=>{const a=JSON.parse(t);H(a).then(l=>{h=!0,document.getElementById("table").innerHTML=l,B(document.getElementById("modeButton").getAttribute("data-mode")),j(),document.querySelectorAll("button").forEach(o=>{o.addEventListener("click",()=>{let n=o.closest("tr");if(o.classList.contains("updateVersion")){let r=document.getElementById("latestScriptVersion").textContent,i=g(n.getAttribute("data-pack-bp-path"))||n.getAttribute("data-pack-path");N(i,n.getAttribute("data-current-version"),r),v(),p(`Updating script for ${n.getAttribute("data-pack-name")}`)}else if(o.classList.contains("compileButton")){let r=n.getAttribute("data-pack-name"),i=n.getAttribute("data-pack-path")||n.getAttribute("data-pack-rp-path");i+="\\pack_icon.png";let s={CleanName:n.getAttribute("data-pack-name"),PackType:n.getAttribute("data-pack-type"),ResoucePackPath:g(n.getAttribute("data-pack-rp-path"))||n.getAttribute("data-pack-path"),BehaviorPackPath:g(n.getAttribute("data-pack-bp-path"))||n.getAttribute("data-pack-path"),ExportPath:"desktop",Format:document.getElementById("saveFormat").value},u=document.getElementById("saveFormat").value;document.getElementById("saveMode").value==="choose"?z({defaultDirectory:"C:/",title:"Select Export Directory"}).then(x=>{s.ExportPath=x,k(JSON.stringify(s)).then(()=>{y(`Finished compiling: 
 ${r}.${u}`,i)})}).catch():k(JSON.stringify(s)).then(()=>{y(`Finished compiling: 
 ${r}.${u}`,i)})}else if(o.classList.contains("normalizePack")){document.getElementById("normalizePanel").style.visibility="visible",T(n.getAttribute("data-pack-rp-path"),n.getAttribute("data-pack-bp-path")).then(i=>{i==="Done"&&(p(`Finished Normalizing: ${n.getAttribute("data-pack-name")}`),document.getElementById("normalizePanel").style.visibility="hidden")}),w("stage:name",i=>{document.getElementById("animatingText").textContent=i});let r=[];w("file:rename",i=>{r.push(i),r.length>10&&r.shift(),document.getElementById("normalizePanelText").innerHTML=r.join("<br><br>")})}})}),document.getElementById("expand").addEventListener("click",()=>{V("https://blockstate.team")})})})}v();window.addEventListener("online",()=>{console.log("Back Online!"),$()});
