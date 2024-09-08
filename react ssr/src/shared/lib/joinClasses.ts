export const joinClasses = (...classes: (string | undefined)[]) => {
    return classes.filter(c => c !== undefined ).join(" ")
}