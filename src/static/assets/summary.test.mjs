import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import vm from "node:vm";

const html = readFileSync(new URL("../../summary.html", import.meta.url), "utf8");
const start = html.indexOf("function updateReviewBanner(periodList)");
const end = html.indexOf("function filterTable()", start);
assert.ok(start >= 0 && end > start, "updateReviewBanner not found");

function element() {
    return {
        classList: { remove() {}, add() {} },
        disabled: false,
        setAttribute() {},
        textContent: "",
    };
}

const elements = {
    reviewQueueBtn: element(),
    reviewQueueText: element(),
    reviewQueueAction: element(),
};
const context = {
    document: { getElementById: (id) => elements[id] },
    needsReview: () => false,
    reviewOnly: false,
    workOrdersLoadFailed: false,
};

vm.runInNewContext(
    `${html.slice(start, end)}
    updateReviewBanner([]);`,
    context,
);
assert.equal(
    elements.reviewQueueText.textContent,
    "Tidak ada work order pada periode ini",
);

vm.runInNewContext("updateReviewBanner([{}]);", context);
assert.equal(
    elements.reviewQueueText.textContent,
    "Semua work order pada periode ini sudah direview",
);

console.log("summary empty-state behavior: ok");
