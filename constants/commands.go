package constants

func Fetch_w() string {
	return `
		Remove-Item $HOME\.config\tran -Recurse -Force
		tran sync fetchx
		Write-Host "Fetched Successfully"
	`
}

func Fetch_ml() string {
	return `
		cd $HOME/.config/tran
		git pull
		echo "Fetched Successfully ✅"
	`
}

func Start_w() string {
	return `
		$username = tran auth get-username
		cd $HOME\.config\tran
		git init
		tran gh-repo create .tran -d "My tran config - $username" --private -y
		git add .
		git commit -m "new .tran repo"
		git branch -M trunk
		git remote add origin https://github.com/$username/.tran
		git push -u origin trunk
		cd $lastDir
	`
}

func Start_ml() string {
	return `
		username=$(tran auth get-username)
		cd $HOME/.config/tran
		git init
		tran gh-repo create .tran -d "My tran config - $username" --private -y
		git add .
		git commit -m "new .tran repo"
		git branch -M trunk
		git remote add origin https://github.com/$username/.tran
		git push -u origin trunk
	`
}

func Push_w() string {
	return `
		$lastDir = pwd
		cd $HOME\.config\tran
		if (Test-Path -path .git) {
			git add .
			git commit -m "new change"
			git push
		}

		cd $lastDir
	`
}

func Push_ml() string {
	return `
		cd $HOME/.config/tran
		git add .
		git commit -m "new tran config"
		git push
	`
}

func Pull_w() string {
	return `
		$lastDir = pwd
		cd $HOME\.config\tran

		git pull

		cd $lastDir
	`
}

func Pull_ml() string {
	return `
		cd $HOME/.config/tran
		git pull
	`
}

func Clone_w() string {
	return `
		$TRANDIR = $HOME\.config\tran

		if (Test-Path -path $TRANDIR) {
			Remove-Item $TRANDIR -Recurse -Force
		} else {
			tran gh-repo clone .tran $TRANDIR
		}
	`
}

func Clone_ml() string {
	return `
		TRANDIR=$HOME/.config/tran

		if [ -d $TRANDIR ]; then
			rm -rf $TRANDIR
		else
			tran gh-repo clone .tran $TRANDIR
		fi
	`
}

func Clone_check_w() string {
	return `
		if (Test-Path -path $HOME\.config\tran) {
			Write-Host "tran repo cloned successfully"
		}
	`
}

func Clone_check_ml() string {
	return `if [ -d $HOME/.config/tran ]; then echo "tran repo cloned successfully ✅"; fi`
}
