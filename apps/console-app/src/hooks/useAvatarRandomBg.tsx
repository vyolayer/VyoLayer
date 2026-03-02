import { useState } from 'react'

const colors = [
    'bg-red-500',
    'bg-green-500',
    'bg-blue-500',
    'bg-yellow-500',
    'bg-purple-500',
    'bg-pink-500',
    'bg-orange-500',
    'bg-indigo-500',
    'bg-teal-500',
    'bg-cyan-500',
    'bg-lime-500',
    'bg-violet-500',
    'bg-fuchsia-500',
    'bg-rose-500',
    'bg-emerald-500',
    'bg-sky-500',
    'bg-amber-500',
]

const AVATAR_BG_COLOR_KEY = 'wl_avatar_bg_color'

/**
 * Returns a random background color for the avatar.
 * Check if already have a color in local storage, if yes then return that color else return a random color and save it in local storage
 *
 * @returns {string} - Random background color
 */
export const useAvatarRandomBg = (): string => {
    const [color] = useState<string>(() => {
        if (typeof window === 'undefined') return colors[0]

        const storedColor = localStorage.getItem(AVATAR_BG_COLOR_KEY)
        if (storedColor && colors.includes(storedColor)) {
            return storedColor
        }

        const newColor = colors[Math.floor(Math.random() * colors.length)]
        localStorage.setItem(AVATAR_BG_COLOR_KEY, newColor)
        return newColor
    })

    return color
}
