import { SelectionMode, SendMessageType } from "../enums.js";
import { networkManager } from "../main.js";
import { SelectChoice } from "../messages/send.js";
import { ChoiceTemplate } from "../modules/choice_template.js";
import { KeyValueTemplate } from "../modules/keyvalue_template.js";

export class SectionMultichoice {
    private itemID: string;
    private switchCount: number = 5;
    private selectedTokenID: string = "";

    private choices: ChoiceTemplate[];

    private section_DOM: HTMLElement;

    private dropdown_DOM: HTMLElement;
    private dropdownSummary_DOM: HTMLElement;
    private dropdownDetails_DOM: HTMLDetailsElement;
    private dropdownTo_DOM: HTMLElement;
    private modalbox_DOM: HTMLElement;
    private modalboxDialog_DOM: HTMLDialogElement;
    private MetadataTable_DOM: HTMLTableElement;
    private modalboxTo_DOM: HTMLElement;


    private manualButton_DOMs: HTMLButtonElement[];
    private confirmButton_DOM: HTMLButtonElement;

    private popupIDSetup(): void {
        const modalid = "modal-" + crypto.randomUUID();
        this.modalbox_DOM.querySelector("dialog")?.setAttribute("id", modalid);
        this.modalbox_DOM.querySelector("button.modal")?.setAttribute("commandfor", modalid);
    }
    constructor(itemID: string, section_DOM: HTMLElement) {
        this.itemID = itemID;
        this.choices = [] as ChoiceTemplate[];
        this.section_DOM = section_DOM;
        this.switchCount = parseInt(this.section_DOM.dataset.switchcount || "5", 10);
        this.dropdown_DOM = this.section_DOM.querySelector("div.select") as HTMLElement;
        this.dropdownSummary_DOM = this.dropdown_DOM.querySelector("summary") as HTMLElement;
        this.dropdownDetails_DOM = this.dropdown_DOM.querySelector("details") as HTMLDetailsElement;
        this.dropdownTo_DOM = this.dropdown_DOM.querySelector("details div") as HTMLElement;
        this.modalbox_DOM = this.section_DOM.querySelector("div.modal") as HTMLElement;
        this.modalboxDialog_DOM = this.section_DOM.querySelector("dialog") as HTMLDialogElement;
        this.MetadataTable_DOM = this.modalbox_DOM.querySelector('table') as HTMLTableElement;
        this.modalboxTo_DOM = this.modalbox_DOM.querySelector("dialog div") as HTMLElement;
        this.manualButton_DOMs = [] as HTMLButtonElement[];
        section_DOM.querySelectorAll("button.manual").forEach(button => {
            this.manualButton_DOMs.push(button as HTMLButtonElement);
        });
        this.confirmButton_DOM = this.dropdown_DOM.querySelector("button.confirm") as HTMLButtonElement;

        this.popupIDSetup();
    }

    private clickManual = (e: Event): void => {
        networkManager.SendPayload(this.itemID, SendMessageType.SwitchToManual);
    }

    private clickConfirm = (e: Event): void => {
        const button = e.currentTarget as HTMLButtonElement;
        button.disabled = true
        const payload: SelectChoice = {
            Choice: this.selectedTokenID,
        };
        networkManager.SendPayload(this.itemID, SendMessageType.SelectChoice, JSON.stringify(payload));
    }

    private callbackModalbox = (choice: ChoiceTemplate): void => {
        this.modalboxDialog_DOM.close();
        const payload: SelectChoice = {
            Choice: choice.TokenID,
        };
        networkManager.SendPayload(this.itemID, SendMessageType.SelectChoice, JSON.stringify(payload));
    };

    private callbackDropdown = (choice: ChoiceTemplate): void => {
        this.dropdownDetails_DOM.open = false;
        this.confirmButton_DOM.disabled = false;
        this.selectedTokenID = choice.TokenID;
        this.dropdownSummary_DOM.replaceChildren(choice.SummaryElement());
    }


    public AddEventListeners(): void {
        this.manualButton_DOMs.forEach(button => {
            button.addEventListener("click", this.clickManual);
        });
        this.confirmButton_DOM.addEventListener("click", this.clickConfirm);
        this.choices.forEach(element => {
            element.AddEventListeners();
        });
    }

    public FoundCountHTML(count: number): void {
        const strong = this.section_DOM.querySelector('h6 strong');
        if (strong) (strong as HTMLElement).innerText = String(count);
        if (count == 0) {
            this.SelectionMode(SelectionMode.Manual);
        } else if (this.switchCount >= count) {
            this.SelectionMode(SelectionMode.Select);
        } else {
            this.SelectionMode(SelectionMode.Modal);
        }
    }

    public SelectionMode(mode: SelectionMode): void {
        this.section_DOM.dataset.viewmode = mode;
    }

    public AddMetadata(keyValues: KeyValueTemplate[]): void {
        keyValues.forEach(keyValue => {
            this.MetadataTable_DOM.appendChild(keyValue.Element());
        });
    }

    public RemoveMetadata(): void {
        this.MetadataTable_DOM.innerHTML = "";
    }

    public AddChoices(choices: ChoiceTemplate[]): void {
        this.choices = choices;
        this.FoundCountHTML(choices.length);
        this.dropdownTo_DOM.innerHTML = "";
        this.modalboxTo_DOM.innerHTML = "";
        const activeModalId = this.modalboxDialog_DOM.getAttribute("id") || "";

        choices.forEach(choice => {
            choice.AddCallbacks(this.callbackDropdown, this.callbackModalbox);
            const choiceBtnNode = choice.ModalElement().querySelector("button");
            if (choiceBtnNode && activeModalId) {
                choiceBtnNode.setAttribute("commandfor", activeModalId);
                choiceBtnNode.setAttribute("command", "show-modal");
            }

            choice.AddEventListeners();
            this.dropdownTo_DOM.appendChild(choice.DropDownElement());
            this.modalboxTo_DOM.appendChild(choice.ModalElement());
        });
    }


    public RemoveChoices(): void {
        this.dropdownTo_DOM.innerHTML = "";
        this.modalboxTo_DOM.innerHTML = "";
        this.choices.length = 0;
    }
}
