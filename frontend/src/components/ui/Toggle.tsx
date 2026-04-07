interface Props {
  checked: boolean
  onChange: (v: boolean) => void
  label?: string
}

export const Toggle = ({ checked, onChange, label }: Props) => (
  <button
    type="button"
    role="switch"
    aria-checked={checked}
    aria-label={label}
    onClick={() => onChange(!checked)}
    className={`relative w-11 h-6 rounded-full cursor-pointer transition-colors duration-200 shrink-0 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/50 focus-visible:ring-offset-2 ${
      checked ? 'bg-blue-600' : 'bg-gray-300'
    }`}
  >
    <span
      className={`absolute top-1 w-4 h-4 bg-white rounded-full shadow transition-all duration-200 ${
        checked ? 'left-6' : 'left-1'
      }`}
    />
  </button>
)
