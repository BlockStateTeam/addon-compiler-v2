import {
    Normalize,
    GetData,
    Notify,
    NotifyText,
    CompilePack,
    OpenDirectoryDialogCall,
    UpdateScriptVersion,
    GetUserSetting,
    SaveUserSetting,
    SelfVersion,
    UninstallApp,
} from "../wailsjs/go/main/App";
import { BrowserOpenURL, EventsOn } from "/wailsjs/runtime/runtime";
const themes = {
    default: {
        "--background-normalize-panel": "rgba(0, 0, 0, 0.5)",
        "--background-color": "#242424",
        "--grid-color": "#ffffff00",
        "--text-color": "#fff",
        "--font-size": "1.5vw",
        "--table-width": "60vw",
        "--table-background": "#171717",
        "--button-background": "#666",
        "--button-hover-background": "#393939",
        "--button-active-background": "#000000",
        "--header-font-size": "4vw",
        "--h2-font-size": "2.5vw",
        "--h3-font-size": "2.75vw",
        "--h3-hover-font-size": "3vw",
        "--h3-hover-color": "#80a2ff",
        "--h4-font-size": "1.5vw",
        "--select-background": "#393939",
        "--select-text-color": "#fff",
        "--default-font": "Arial, sans-serif",
        "--h1-font": "Arial, sans-serif",
    },
    light: {
        "--background-normalize-panel": "rgba(255, 255, 255, 0.5)",
        "--background-color": "#ffffff",
        "--grid-color": "#00000000",
        "--text-color": "#242424",
        "--font-size": "1.5vw",
        "--table-width": "60vw",
        "--table-background": "#d5d5d5",
        "--button-background": "#cacaca",
        "--button-hover-background": "#ffffff",
        "--button-active-background": "#c6c6c6",
        "--header-font-size": "4vw",
        "--h2-font-size": "2.5vw",
        "--h3-font-size": "2.75vw",
        "--h3-hover-font-size": "3vw",
        "--h3-hover-color": "#1654ff",
        "--h4-font-size": "1.5vw",
        "--select-background": "#393939",
        "--select-text-color": "#fff",
        "--default-font": "Arial, sans-serif",
        "--h1-font": "Arial, sans-serif",
    },
    minecraft: {
        "--background-normalize-panel": "rgba(0, 0, 0, 0.5)",
        "--background-color": "#242424",
        "--grid-color": "#ffffff00",
        "--text-color": "#fff",
        "--font-size": "1.3vw",
        "--table-width": "60vw",
        "--table-background": "#171717",
        "--button-background": "#666",
        "--button-hover-background": "#393939",
        "--button-active-background": "#c9c9c9",
        "--header-font-size": "4vw",
        "--h2-font-size": "2.5vw",
        "--h3-font-size": "2.75vw",
        "--h3-hover-font-size": "3vw",
        "--h3-hover-color": "#80a2ff",
        "--h4-font-size": "1.3vw",
        "--select-background": "#393939",
        "--select-text-color": "#fff",
        "--default-font": "Minecraft",
        "--h1-font": "MinecraftTen",
    },
    grid: {
        "--background-normalize-panel": "rgba(0, 0, 0, 0.5)",
        "--background-color": "#242424",
        "--grid-color": "#ffffff10",
        "--text-color": "#fff",
        "--font-size": "1.5vw",
        "--table-width": "60vw",
        "--table-background": "#171717",
        "--button-background": "#666",
        "--button-hover-background": "#393939",
        "--button-active-background": "#000000",
        "--header-font-size": "4vw",
        "--h2-font-size": "2.5vw",
        "--h3-font-size": "2.75vw",
        "--h3-hover-font-size": "3vw",
        "--h3-hover-color": "#80a2ff",
        "--h4-font-size": "1.5vw",
        "--select-background": "#393939",
        "--select-text-color": "#fff",
        "--default-font": "Arial, sans-serif",
        "--h1-font": "Arial, sans-serif",
    },
};
let userSetting;

function applyTheme(themeVariables) {
    const root = document.documentElement;
    for (const [name, value] of Object.entries(themeVariables)) {
        root.style.setProperty(name, value);
    }
}
SelfVersion().then((version) => {
    document.getElementById("selfVersion").textContent = "v" + version;
});
GetUserSetting().then((data) => {
    userSetting = JSON.parse(data);
    document.getElementById("saveFormat").value =
        userSetting.format || "mcaddon";
    document.getElementById("saveMode").value =
        userSetting.exportLocation || "desktop";
    document.getElementById("uiselect").value = userSetting.theme || "default";
    applyTheme(themes[userSetting.theme]);
});

function updateUserSetting(type, value) {
    userSetting[type] = value;
    SaveUserSetting(JSON.stringify(userSetting));
}
function toggleModeState(event) {
    const button = event.currentTarget;
    if (button.textContent === "World Mode") {
        button.textContent = "Add-On Mode";
        document.querySelectorAll(".isWorld").forEach((element) => {
            element.style.display = "none";
        });
        document.querySelectorAll(".isNotWorld").forEach((element) => {
            element.style.display = "block"; // Or 'flex', 'grid' etc.
        });
    } else {
        button.textContent = "World Mode";
        document.querySelectorAll(".isWorld").forEach((element) => {
            element.style.display = "block";
        });
        document.querySelectorAll(".isNotWorld").forEach((element) => {
            element.style.display = "none";
        });
    }
}
document.getElementById("modeButton").addEventListener("click", (event) => {
    toggleModeState(event);
});
document.getElementById("saveFormat").addEventListener("change", (select) => {
    updateUserSetting("format", select.target.value);
});
document.getElementById("saveMode").addEventListener("change", (select) => {
    updateUserSetting("exportLocation", select.target.value);
});
document.getElementById("uiselect").addEventListener("change", (select) => {
    let value = select.target.value;
    updateUserSetting("theme", value);
    applyTheme(themes[value]);
});
async function generateHTMLTables(inputArray) {
    const tablesHTML = [];
    for (let i = 0; i <= 3; i++) {
        let array = inputArray[i];
        const isWorldClass = i === 3 ? "isWorld" : "isNotWorld";
        let tableHTML = `
            <h2 class="${isWorldClass}">
                ${["Resource Pack", "Add-On Pack", "Behavior Pack", "Worlds", "Warning"][i]}
            </h2>
            <table class="${isWorldClass}">
                <tbody>`;
        if (array) {
            for (let j = 0; j < array.length; j++) {
                const pack = array[j];
                const cleanName = pack.CleanName;
                const scriptState =
                    pack?.ScriptState || pack?.BehaviorPack?.ScriptState;
                let scriptCurrentVersion = scriptState
                    ? scriptState
                    : "null-script";
                switch (i) {
                    case 1: // Add-On Pack
                        tableHTML += `
                            <tr data-current-version="${scriptCurrentVersion}" data-pack-name="${cleanName}" data-pack-type="addOnPack" data-pack-path="${base64EncodeUnicode(JSON.stringify([pack.ResourcePack.Path, pack.BehaviorPack.Path]))}" data-pack-signatures="${pack.IsSignatures}">
                                <td class="noneData"><img class="projectImage" src="/${base64EncodeUnicode(pack.BehaviorPack.Icon)}.png"></td>
                                <td id="textCell" class="averageTextCell">${cleanName}
                                <div class="hiddenAddonAttributes">
                                    <div>
                                        <span class="packNote">Behavior Pack:</span><br>
                                        <span><img class="projectImage" src="/${base64EncodeUnicode(pack.BehaviorPack.Icon)}.png"> ${pack.BehaviorPack.CleanName}</span>
                                    </div>
                                    <br>
                                    <div>
                                        <span class="packNote">Resource Pack:</span><br>
                                        <span><img class="projectImage" src="/${base64EncodeUnicode(pack.ResourcePack.Icon)}.png"> ${pack.ResourcePack.CleanName}</span>
                                    </div>
                                </div>
                        `;
                        break;
                    case 3: // World
                        //console.log(pack);
                        const requiredRP = JSON.parse(
                            base64Decode(pack.RequiredRP),
                        );
                        const requiredBP = JSON.parse(
                            base64Decode(pack.RequiredBP),
                        );
                        console.log("requiredBP: ", requiredBP);
                        console.log("requiredRP: ", requiredRP);
                        tableHTML += `
                            <tr data-pack-name="${cleanName}" data-pack-type="world" data-pack-required-rp="${pack.RequiredRP}" data-pack-required-bp="${pack.RequiredBP}" data-pack-path="${base64EncodeUnicode(JSON.stringify([pack.Path]))}">
                                <td class="noneData"><img class="worldImage" src="/${base64EncodeUnicode(pack.Icon)}.png"></td>
                                <td id="textCell" class="averageTextCell">${cleanName}
                                <div class="hiddenAddonAttributes">
                                    <div>
                        `;
                        if (requiredBP.length > 0) {
                            tableHTML += `<span class="packNote">Behavior Pack:</span><br>`;
                            requiredBP.forEach((bp) => {
                                tableHTML += `
                                    <span><img class="projectImage" src="/${base64EncodeUnicode(bp.Icon)}.png"> ${bp.CleanName}</span>
                                `;
                            });
                        }
                        tableHTML += `</div><br><div>`;
                        if (requiredRP.length > 0) {
                            tableHTML += `<span class="packNote">Resource Pack:</span><br>`;
                            requiredRP.forEach((rp) => {
                                tableHTML += `
                                    <span><img class="projectImage" src="/${base64EncodeUnicode(rp.Icon)}.png"> ${rp.CleanName}</span>
                                `;
                            });
                        }
                        tableHTML += `</div></div>`;
                        break;
                    default: // Resource Pack & Behavior Pack
                        const type = i === 0 ? "resourcePack" : "behaviorPack";
                        tableHTML += `
                            <tr data-current-version="${scriptCurrentVersion}" data-pack-name="${cleanName}" data-pack-type="${type}" data-pack-path="${base64EncodeUnicode(JSON.stringify([pack.Path]))}" data-pack-signatures="${pack.IsSignatures}">
                                <td class="noneData"><img class="projectImage" src="/${base64EncodeUnicode(pack.Icon)}.png"></td>
                                <td id="textCell" class="averageTextCell">${cleanName}
                        `;
                }
                tableHTML += `
                        </td>
                        <td id="buttonCell1"><button class="updateVersion" style="visibility:hidden;">Update API</button></td>
                        <td id="buttonCell2"><button class="compileButton">Compile</button></td>
                        <td id="buttonCell3" class="normalizeButton"><button class="normalizePack" style="visibility:hidden;">Normalize Pack</button></td>
                    </tr>
                `;
            }
        }
        tableHTML += `</tbody></table>`;
        tablesHTML.push(tableHTML);
    }
    return tablesHTML.join("");
}
function reload() {
    GetData().then((data) => {
        const result = JSON.parse(data);
        generateHTMLTables(result).then((tables) => {
            document.getElementById("table").innerHTML = tables;
            document.getElementById("modeButton").textContent = "Add-On Mode";
            document.querySelectorAll("button").forEach((button) => {
                button.addEventListener("click", () => {
                    const dataElement = button.closest("tr");
                    const packIcon = JSON.parse(
                        base64Decode(
                            dataElement.getAttribute("data-pack-path"),
                        ),
                    )[0];
                    if (button.classList.contains("updateVersion")) {
                        UpdateScriptVersion(
                            dataElement.getAttribute("data-pack-path"),
                            dataElement.getAttribute("data-current-version"),
                        ).then((result) => {
                            Notify(result, packIcon);
                        });
                        reload();
                    } else if (button.classList.contains("compileButton")) {
                        let packData = {
                            CleanName:
                                dataElement.getAttribute("data-pack-name"),
                            PackType:
                                dataElement.getAttribute("data-pack-type"),
                            PackPath:
                                dataElement.getAttribute("data-pack-path"),
                            ExportPath: "desktop",
                            RequiredRP: dataElement.getAttribute(
                                "data-pack-required-rp",
                            ),
                            RequiredBP: dataElement.getAttribute(
                                "data-pack-required-bp",
                            ),
                            Format: document.getElementById("saveFormat").value,
                        };
                        const format =
                            document.getElementById("saveFormat").value;
                        const packDataJson = JSON.stringify(packData);
                        if (
                            document.getElementById("saveMode").value ===
                            "choose"
                        ) {
                            OpenDirectoryDialogCall()
                                .then((path) => {
                                    packData.ExportPath = path;
                                    CompilePack(packDataJson).then((result) => {
                                        Notify(result, packIcon);
                                    });
                                })
                                .catch();
                        } else {
                            CompilePack(packDataJson).then((result) => {
                                Notify(result, packIcon);
                            });
                        }
                    } else if (button.classList.contains("normalizePack")) {
                        document.getElementById(
                            "normalizePanel",
                        ).style.visibility = "visible";
                        Normalize(
                            dataElement.getAttribute("data-pack-rp-path"),
                            dataElement.getAttribute("data-pack-bp-path"),
                        ).then((result) => {
                            if (result === "Done") {
                                NotifyText(
                                    `Finished Normalizing: ${dataElement.getAttribute("data-pack-name")}`,
                                );
                                document.getElementById(
                                    "normalizePanel",
                                ).style.visibility = "hidden";
                            }
                        });
                        EventsOn("stage:name", (stage) => {
                            document.getElementById(
                                "animatingText",
                            ).textContent = stage;
                        });
                        let dataLines = [];
                        EventsOn("file:rename", (data) => {
                            dataLines.push(data);
                            if (dataLines.length > 10) {
                                dataLines.shift();
                            }
                            document.getElementById(
                                "normalizePanelText",
                            ).innerHTML = dataLines.join("<br><br>");
                        });
                    }
                });
            });
            document.querySelectorAll(".updateVersion").forEach((button) => {
                // button is inside a td, which is inside a tr
                let dataElement = button.closest("tr");
                let currentVersion = dataElement.getAttribute(
                    "data-current-version",
                );
                if (currentVersion?.endsWith("-beta")) {
                    button.style.visibility = "visible";
                    button.setAttribute(
                        "data-pack-path",
                        dataElement.getAttribute("data-pack-path"),
                    );
                    button.setAttribute("data-current-version", currentVersion);
                }
            });
            document.getElementById("expand").addEventListener("click", () => {
                BrowserOpenURL("https://blockstate.team");
            });
        });
    });
}
reload();
document.getElementById("reload").addEventListener("click", () => {
    reload();
});
function base64EncodeUnicode(str) {
    const utf8String = encodeURIComponent(str).replace(
        /%([0-9A-F]{2})/g,
        function (match, p1) {
            return String.fromCharCode("0x" + p1);
        },
    );
    return btoa(utf8String);
}
function base64Decode(base64String) {
    const binaryString = atob(base64String);
    return decodeURIComponent(
        binaryString
            .split("")
            .map(function (c) {
                return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
            })
            .join(""),
    );
}
document.getElementById("uninstall").addEventListener("click", () => {
    document.getElementById("uninstallPanel").style.visibility = "visible";
});
document.getElementById("confirmUninstall").addEventListener("click", () => {
    UninstallApp().then((result) => {
        if (result != "") {
            NotifyText(result);
        }
    });
});
document.getElementById("cancelUninstall").addEventListener("click", () => {
    document.getElementById("uninstallPanel").style.visibility = "hidden";
});
