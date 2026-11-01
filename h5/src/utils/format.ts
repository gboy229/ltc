/**
 * 金额展示：固定小数位，不使用千分位分隔符。
 * 例如：5000.00、1234.5 -> 1234.50、0.8 -> 0.80
 */
export function formatAmount(value: number, digits = 4): string {
  if (!Number.isFinite(value)) {
    return (0).toFixed(digits)
  }
  return value.toFixed(digits)
}
