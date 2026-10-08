export class ListItemTemplate {
    private static template_DOM: HTMLTemplateElement;

    public static SetTemplate(template: HTMLTemplateElement): void {
        ListItemTemplate.template_DOM = template;
    }

    private value: string;
    private icon: string;

    constructor(value: string, icon: string) {
        this.value = value as string;
        this.icon = icon as string;
    }

    public Element(): HTMLLIElement {
        const li_DOM = document.importNode(ListItemTemplate.template_DOM.content, true).querySelector("li") as HTMLLIElement;
        li_DOM.innerHTML = this.value;
        li_DOM.style.setProperty('--bullet-icon', `"${this.icon} "`);
        return li_DOM as HTMLLIElement;
    }
}
