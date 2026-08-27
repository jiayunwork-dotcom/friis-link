"use strict";
const $ = (id) => document.getElementById(id);
function showError(msg) { $("err-text").textContent = msg; $("err-panel").hidden = false; }
function hideError() { $("err-panel").hidden = true; }
function fill(rows) {
  const tb = $("out-table"); tb.innerHTML = "";
  for (const [k,v] of rows) { const tr = document.createElement("tr"); tr.innerHTML = "<td>"+k+"</td><td>"+v+"</td>"; tb.appendChild(tr); }
}
async function loadExample() {
  hideError();
  try {
    const resp = await fetch("/api/example"); const data = await resp.json();
    if (!resp.ok) throw new Error(data.error || "示例失败");
    $("body").value = JSON.stringify(data, null, 2);
    $("hint").textContent = "已加载 S 波段 10 km 算例。";
  } catch (e) { showError(String(e)); }
}
async function runBudget() {
  hideError();
  try {
    const resp = await fetch("/api/budget", { method:"POST", headers:{"Content-Type":"application/json"}, body: $("body").value });
    const data = await resp.json();
    if (!resp.ok) throw new Error(data.error || "HTTP "+resp.status);
    $("out-panel").hidden = false;
    fill([
      ["波长 m", data.lambda_m.toPrecision(6)],
      ["FSPL dB", data.fspl_db.toPrecision(6)],
      ["EIRP dBm", data.eirp_dbm.toPrecision(6)],
      ["Pr dBm", data.pr_dbm.toPrecision(6)]
    ]);
  } catch (e) { showError(String(e)); }
}
$("btn-example").addEventListener("click", loadExample);
$("btn-run").addEventListener("click", runBudget);
loadExample();
