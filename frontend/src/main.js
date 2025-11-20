import {
    Normalize,
    GetData,
    GetImage,
    Notify,
    NotifyText,
    CompilePack,
    OpenDirectoryDialog,
    UpdateScriptVersion,
    GetUserSetting,
    SaveUserSetting,
    SelfVersion,
} from "../wailsjs/go/main/App";
import {
    BrowserOpenURL,
    WindowReload,
    EventsOn,
} from "/wailsjs/runtime/runtime";
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
let freeTable = true;
let poi = ["BUTTON", "SELECT", "SPAN", "TD"];
let contextMenu = document.getElementById("contextMenu");
function insertSpaceBeforeCapitals(str) {
    str = str.charAt(0).toUpperCase() + str.slice(1);
    return str.replace(/([a-z])([A-Z])/g, "$1 $2");
}
document.addEventListener("mouseover", (event) => {
    let target = event.target;
    if (poi.includes(target.tagName)) {
        switch (target.tagName) {
            case "TD":
                let parent = target.closest("tr");
                let scriptApiOfPack =
                    parent.getAttribute("data-current-version") ===
                    "null-script"
                        ? ""
                        : " with BetaAPI version " +
                          parent.getAttribute("data-current-version");
                contextMenu.innerHTML = `<strong>${insertSpaceBeforeCapitals(parent.getAttribute("data-pack-type"))}</strong> ${scriptApiOfPack}`;
                break;
            case "BUTTON":
                let textContent = target.textContent;
                if (textContent === "Update API") {
                    contextMenu.innerHTML =
                        "Click to <strong>update</strong> the Script API version for this pack";
                } else if (textContent === "Compile") {
                    let parent = target.closest("tr");
                    let choice =
                        document.getElementById("saveMode").value === "desktop"
                            ? "to Desktop"
                            : "";
                    contextMenu.innerHTML = `Click to compile <strong>${parent.getAttribute("data-pack-name")}</strong> as ${document.getElementById("saveFormat").value.toUpperCase()} format ${choice}`;
                } else if (textContent === "Dev Mode") {
                    contextMenu.innerHTML =
                        "Click to toggle between <strong>Development</strong> and <strong>Release</strong> mode";
                }
                break;
            case "SELECT":
                if (target.id === "saveFormat") {
                    contextMenu.innerHTML = `Change export <strong>format</strong> for compiled packs`;
                } else if (target.id === "saveMode") {
                    contextMenu.innerHTML = `Change export <strong>location</strong> for compiled packs`;
                } else {
                    contextMenu.innerHTML =
                        "Change <strong>theme</strong> for the app";
                }
                break;
            case "SPAN":
                if (target.id === "latestScriptVersion") {
                    contextMenu.innerHTML =
                        "Script API version refer to <strong>BetaAPI</strong> experimental features.<br><br> Version: " +
                        target.textContent;
                }
                break;
        }
    } else {
        contextMenu.innerHTML =
            "Hover above buttons or selection <br> to get more info <br><br> Click to hide this pannel";
    }
});
document.addEventListener("mousemove", (event) => {
    let offsetTotalWidth = window.screen.width - (event.clientX + 425);
    let offsetTotalHeight =
        window.screen.height -
        (event.clientY +
            document.getElementById("contextMenu").offsetHeight +
            75);
    if (offsetTotalWidth < 0) {
        offsetTotalWidth = event.clientX - 425;
    } else {
        offsetTotalWidth = event.clientX + 5;
    }
    if (offsetTotalHeight < 0) {
        offsetTotalHeight =
            event.clientY -
            document.getElementById("contextMenu").offsetHeight -
            5;
    } else {
        offsetTotalHeight = event.clientY + 5;
    }
    contextMenu.style.top = `${offsetTotalHeight}px`;
    contextMenu.style.left = `${offsetTotalWidth}px`;
});
let contextMenuVisible = false;
document.addEventListener("contextmenu", (event) => {
    event.preventDefault();
    if (contextMenuVisible) {
        document.getElementById("contextMenu").style.opacity = "0";
        contextMenuVisible = false;
        return;
    }
    contextMenuVisible = true;
    let offsetTotal = window.screen.width - (event.clientX + 425);
    if (offsetTotal < 0) {
        offsetTotal = event.clientX - 420;
    } else {
        offsetTotal = event.clientX;
    }
    contextMenu.style.display = "block";
    contextMenu.style.opacity = "1";
    document.addEventListener("click", () => {
        contextMenuVisible = false;
        document.getElementById("contextMenu").style.opacity = "0";
    });
});

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
    document.getElementById("modeButton").getAttribute("data-mode") ===
    userSetting.mode
        ? null
        : toggleModeState();
    applyTheme(themes[userSetting.theme]);
});
function updateTableName(mode) {
    let suffix = "";
    if (mode === "development") {
        suffix = "Development ";
    }
    if (document.getElementById("resourcePack")) {
        document.getElementById("resourcePack").textContent =
            suffix + "Resource Packs";
        document.getElementById("addOnPack").textContent =
            suffix + "Add-On Packs";
        document.getElementById("behaviorPack").textContent =
            suffix + "Behavior Packs";
    }
}
function updateUserSetting(type, value) {
    userSetting[type] = value;
    SaveUserSetting(JSON.stringify(userSetting));
}
function toggleModeState() {
    if (freeTable) {
        console.log("Toggling Mode");
        let button = document.getElementById("modeButton");
        if (button.getAttribute("data-mode") === "release") {
            button.setAttribute("data-mode", "development");
            button.textContent = "Dev Mode";
        } else {
            button.setAttribute("data-mode", "release");
            button.textContent = "Release Mode";
        }
        try {
            updateUserSetting("mode", button.getAttribute("data-mode"));
            reload();
        } catch (error) {
            console.error(error);
        }
        updateTableName(button.getAttribute("data-mode"));
    } else {
        setTimeout(toggleModeState, 100);
    }
}
document.getElementById("modeButton").addEventListener("click", (e) => {
    toggleModeState();
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
    let tableName = ["resourcePack", "addOnPack", "behaviorPack"];
    const tablesHTML = [];
    for (let i = 0; i < inputArray.length; i++) {
        let array = inputArray[i];
        let tableHTML = `<h2 id="${tableName[i]}">Error</h2><table><tbody>`;
        if (array) {
            for (let j = 0; j < array.length; j++) {
                let pack = array[j];
                let scriptCurrentVersion = "null-script";
                if (pack.ScriptState) scriptCurrentVersion = pack.ScriptState;
                let projectImageBase = await GetImage(
                    pack.Path || pack.ResourcePath,
                );
                tableHTML += `
                    <tr data-current-version="${scriptCurrentVersion}" data-pack-name="${pack.CleanName}" data-pack-type="${tableName[i]}" data-pack-path="${pack.Path}" data-pack-rp-path="${pack.ResourcePath}" data-pack-bp-path="${pack.BehaviorPath}" data-pack-signatures="${pack.IsSignatures}">
                        <td class="noneData"><img class="projectImage" src="${projectImageBase}"></td>
                        <td id="textCell" class="averageTextCell">${pack.CleanName}</td>
                        <td id="buttonCell1" class="updateScriptVersionButton"><button class="updateVersion" style="visibility:hidden;">Update API</button></td>
                        <td id="buttonCell2"><button class="compileButton">Compile</button></td>
                        <td id="buttonCell3" class="normalizeButton"><button class="normalizePack" style="visibility:hidden;">Normalize Pack</button></td>
                    </tr>`;
            }
        }
        tableHTML += `</tbody></table>`;
        tablesHTML.push(tableHTML);
    }
    return tablesHTML.join("");
}
function updateCheck() {
    let smallerStrings = [];
    fetch("https://registry.npmjs.org/@minecraft/server")
        .then((response) => {
            if (!response.ok) {
                NotifyText(
                    "Failed to check for Script API version. Please check your internet connection.",
                );
                throw new Error("Network response was not ok");
            }
            return response.json();
        })
        .then((data) => {
            try {
                const numString = Math.max(
                    ...Object.keys(data.time)
                        .filter((key) => /^\d+\.\d+\.\d+/.test(key))
                        .filter((item) => item.endsWith("-stable"))
                        .map((str) => {
                            const match = str.match(/^(\d+\.\d+\.\d+)/);
                            return match
                                ? match[1]
                                      .split(".")
                                      .reduce(
                                          (acc, num, index) =>
                                              acc +
                                              num * Math.pow(100, 2 - index),
                                          0,
                                      )
                                : null;
                        }),
                ).toString();
                let i = numString.length;
                for (; i > 0; i -= 2) {
                    smallerStrings.unshift(
                        numString.slice(Math.max(i - 2, 0), i),
                    );
                }
                let latestVersion =
                    smallerStrings.map((item) => Number(item)).join(".") +
                    "-beta";
                document.getElementById("latestScriptVersion").innerText =
                    latestVersion;
                document.getElementById(
                    "latestScriptContainer",
                ).style.visibility = "visible";
                [...document.getElementsByClassName("normalizeButton")].forEach(
                    (button) => {
                        let boolValue =
                            button
                                .closest("tr")
                                .getAttribute("data-pack-signatures") === "true"
                                ? true
                                : false;
                        if (boolValue) {
                            console.log(boolValue);
                            button.querySelector("button").style.visibility =
                                "visible";
                        }
                    },
                );
                [
                    ...document.getElementsByClassName(
                        "updateScriptVersionButton",
                    ),
                ].forEach((button) => {
                    let currentVersion = button
                        .closest("tr")
                        .getAttribute("data-current-version");
                    if (
                        currentVersion !== "null-script" &&
                        currentVersion.endsWith("-beta") &&
                        currentVersion !== latestVersion
                    ) {
                        button.querySelector("button").style.visibility =
                            "visible";
                    }
                });
            } catch (error) {
                console.error("Error: " + error);
            }
        })
        .catch((error) => {
            console.warn(error);
        });
}
function checkUndefined(value) {
    return value === "undefined" ? undefined : value;
}
function reload() {
    freeTable = false;
    let isDevMode =
        document.getElementById("modeButton").getAttribute("data-mode") ===
        "development";
    console.log(isDevMode);
    GetData(isDevMode).then((data) => {
        const result = JSON.parse(data);
        generateHTMLTables(result).then((tables) => {
            freeTable = true;
            document.getElementById("table").innerHTML = tables;
            updateTableName(
                document.getElementById("modeButton").getAttribute("data-mode"),
            );
            updateCheck();
            document.querySelectorAll("button").forEach((button) => {
                button.addEventListener("click", () => {
                    let dataElement = button.closest("tr");
                    if (button.classList.contains("updateVersion")) {
                        let latestScriptVersion = document.getElementById(
                            "latestScriptVersion",
                        ).textContent;
                        let behaviorPath =
                            checkUndefined(
                                dataElement.getAttribute("data-pack-bp-path"),
                            ) || dataElement.getAttribute("data-pack-path");
                        UpdateScriptVersion(
                            behaviorPath,
                            dataElement.getAttribute("data-current-version"),
                        );
                        reload();
                        NotifyText(
                            `Updating script for ${dataElement.getAttribute("data-pack-name")}`,
                        );
                    } else if (button.classList.contains("compileButton")) {
                        let packName =
                            dataElement.getAttribute("data-pack-name");
                        let packIcon =
                            dataElement.getAttribute("data-pack-path") ||
                            dataElement.getAttribute("data-pack-rp-path");
                        packIcon += "\\pack_icon.png";
                        let packData = {
                            CleanName:
                                dataElement.getAttribute("data-pack-name"),
                            PackType:
                                dataElement.getAttribute("data-pack-type"),
                            ResoucePackPath:
                                checkUndefined(
                                    dataElement.getAttribute(
                                        "data-pack-rp-path",
                                    ),
                                ) || dataElement.getAttribute("data-pack-path"),
                            BehaviorPackPath:
                                checkUndefined(
                                    dataElement.getAttribute(
                                        "data-pack-bp-path",
                                    ),
                                ) || dataElement.getAttribute("data-pack-path"),
                            ExportPath: "desktop",
                            Format: document.getElementById("saveFormat").value,
                        };
                        let format =
                            document.getElementById("saveFormat").value;
                        if (
                            document.getElementById("saveMode").value ===
                            "choose"
                        ) {
                            OpenDirectoryDialog()
                                .then((path) => {
                                    packData.ExportPath = path;
                                    CompilePack(JSON.stringify(packData)).then(
                                        () => {
                                            Notify(
                                                `Finished compiling: \n ${packName}.${format}`,
                                                packIcon,
                                            );
                                        },
                                    );
                                })
                                .catch();
                        } else {
                            CompilePack(JSON.stringify(packData)).then(() => {
                                Notify(
                                    `Finished compiling: \n ${packName}.${format}`,
                                    packIcon,
                                );
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
            document.getElementById("expand").addEventListener("click", () => {
                BrowserOpenURL("https://blockstate.team");
            });
        });
    });
}
reload();
window.addEventListener("online", () => {
    console.log("Back Online!");
    WindowReload();
});
