export function GetTime() {
    const currentDate = new Date();
    const datetime = new Intl.DateTimeFormat("en-US", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit"
    });
    return datetime.format(currentDate);
}