export const enum Status {
    Detecting,
    Multichoice,
    Manual,
    Converting,
    Verifying,
    Verified,
    Failed,
    Apidown,
}

export const StatusFromInt: string[] = [
    "detecting",
    "multichoice",
    "manual",
    "converting",
    "verifying",
    "verified",
    "failed",
    "apidown"
];

export const enum SelectionMode {
    Select = "select",
    Modal = "modal",
    Manual = "manual"
}

export const enum RecieveMessageType {
    CreateItem,
    RemoveItem,
    AcceptToken,
    AddMetadata,
    RemoveMetadata,
    ChangeStatus,
    DetectingAddLog,
    DetectingRemoveLogs,
    AddChoices,
    RemoveChoices,
}

export const enum SendMessageType {
    Ping,
    DetectingStuck,
    SubmitToken,
    FailedRetry,
    FailedDelete,
    SwitchToManual,
    SelectChoice,
    ManualSearch,
    Rescan,
}
