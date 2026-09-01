import fs from 'fs'
import path from 'path'

try {
	const wailsPath = path.resolve('wails.json')
	const wailsData = JSON.parse(fs.readFileSync(wailsPath, 'utf8'))
	const rawVersion = wailsData.version?.trim()

	if (!rawVersion) {
		process.exit(1)
	}

	const goVersion =
		rawVersion.startsWith('v') || rawVersion.startsWith('V') ? rawVersion : `v${rawVersion}`

	const goFilePath = path.resolve('internal/updates/updates.go')
	let goFileContent = fs.readFileSync(goFilePath, 'utf8')
	const goRegex = /currentVersion\s*=\s*"[^"]+"/

	if (!goRegex.test(goFileContent)) {
		process.exit(1)
	}

	goFileContent = goFileContent.replace(goRegex, `currentVersion  = "${goVersion}"`)
	fs.writeFileSync(goFilePath, goFileContent, 'utf8')

	const frontPath = path.resolve('frontend/package.json')
	const frontData = JSON.parse(fs.readFileSync(frontPath, 'utf8'))
	frontData.version = rawVersion
	fs.writeFileSync(frontPath, JSON.stringify(frontData, null, '\t'), 'utf8')

	const lockPath = path.resolve('frontend/package-lock.json')
	if (fs.existsSync(lockPath)) {
		const lockData = JSON.parse(fs.readFileSync(lockPath, 'utf8'))

		if (lockData.version !== undefined) lockData.version = rawVersion
		if (lockData.packages && lockData.packages[''] && lockData.packages[''].version !== undefined) {
			lockData.packages[''].version = rawVersion
		}

		fs.writeFileSync(lockPath, JSON.stringify(lockData, null, '\t'), 'utf8')
	}
} catch (error) {
	process.exit(1)
}
