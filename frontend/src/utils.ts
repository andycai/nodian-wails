export function truncateName(name: string, maxLength: number): string {
    if (name.length <= maxLength) return name;
    const isChinese = /[\u4e00-\u9fa5]/.test(name);
    const limit = isChinese ? 6 : 10;
    const extension = name.split('.').pop();
    const nameWithoutExt = name.substring(0, name.length - extension.length - 1);
    return nameWithoutExt.slice(0, limit) + '...' + (extension ? `.${extension}` : '');
}

// function truncateName(name: string, maxLength: number): string {
//     if (name.length <= maxLength) return name;
//     const isChinese = /[\u4e00-\u9fa5]/.test(name);
//     const limit = isChinese ? 8 : 12;
//     const extension = name.split(".").pop();
//     const nameWithoutExt = name.substring(
//         0,
//         name.length - extension.length - 1,
//     );
//     return (
//         nameWithoutExt.slice(0, limit) +
//         "..." +
//         (extension ? `.${extension}` : "")
//     );
// }