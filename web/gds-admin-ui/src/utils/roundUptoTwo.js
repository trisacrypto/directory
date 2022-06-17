
export default function roundUpToTwo(num) {
    if (isNaN(num)) {
        console.error(`"${num}" is not a number`)
        return
    }

    return Math.round((num + Number.EPSILON) * 100) / 100
}
