// Code generated by build.mjs; DO NOT EDIT.

package assets

import _ "embed"

//go:embed "assets_disclaimer.txt"
var Disclaimer string

// Public holds the path to the public route 
const Public = "/⛰/"

// AssetsDefault contains assets for the 'Default' entrypoint.
var AssetsDefault = Assets{
	Scripts: `<script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule defer></script><script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule defer></script>`,
	Styles:  `<link rel="stylesheet" href="/⛰/Default.8912112e.css"><link rel="stylesheet" href="/⛰/Default.d8ae99b3.css">`,	
}

// AssetsUser contains assets for the 'User' entrypoint.
var AssetsUser = Assets{
	Scripts: `<script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule defer></script><script type="module" src="/⛰/User.23a71b44.js"></script><script src="/⛰/User.869c9a74.js" nomodule defer></script>`,
	Styles:  `<link rel="stylesheet" href="/⛰/Default.8912112e.css"><link rel="stylesheet" href="/⛰/User.68febbf8.css"><link rel="stylesheet" href="/⛰/User.09b09c46.css">`,	
}

// AssetsAdmin contains assets for the 'Admin' entrypoint.
var AssetsAdmin = Assets{
	Scripts: `<script nomodule defer src="/⛰/User.869c9a74.js"></script><script type="module" src="/⛰/User.23a71b44.js"></script><script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule defer></script><script type="module" src="/⛰/Admin.c0a122d2.js"></script><script src="/⛰/Admin.b992cf94.js" nomodule defer></script>`,
	Styles:  `<link rel="stylesheet" href="/⛰/Default.8912112e.css"><link rel="stylesheet" href="/⛰/User.09b09c46.css"><link rel="stylesheet" href="/⛰/Admin.9d294a13.css"><link rel="stylesheet" href="/⛰/User.68febbf8.css"><link rel="stylesheet" href="/⛰/Admin.db3e959a.css">`,	
}

// AssetsAdminProvision contains assets for the 'AdminProvision' entrypoint.
var AssetsAdminProvision = Assets{
	Scripts: `<script nomodule defer src="/⛰/User.869c9a74.js"></script><script nomodule defer src="/⛰/Admin.b992cf94.js"></script><script type="module" src="/⛰/User.23a71b44.js"></script><script type="module" src="/⛰/Admin.c0a122d2.js"></script><script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule defer></script><script type="module" src="/⛰/AdminProvision.65c71f50.js"></script><script src="/⛰/AdminProvision.d01fac9b.js" nomodule defer></script>`,
	Styles:  `<link rel="stylesheet" href="/⛰/Default.8912112e.css"><link rel="stylesheet" href="/⛰/Admin.db3e959a.css"><link rel="stylesheet" href="/⛰/User.68febbf8.css"><link rel="stylesheet" href="/⛰/User.09b09c46.css"><link rel="stylesheet" href="/⛰/Admin.9d294a13.css"><link rel="stylesheet" href="/⛰/AdminProvision.38d394c2.css">`,	
}

// AssetsAdminRebuild contains assets for the 'AdminRebuild' entrypoint.
var AssetsAdminRebuild = Assets{
	Scripts: `<script nomodule defer src="/⛰/User.869c9a74.js"></script><script nomodule defer src="/⛰/Admin.b992cf94.js"></script><script type="module" src="/⛰/User.23a71b44.js"></script><script type="module" src="/⛰/Admin.c0a122d2.js"></script><script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule defer></script><script type="module" src="/⛰/AdminRebuild.d7ba1ab0.js"></script><script src="/⛰/AdminRebuild.0be08ae3.js" nomodule defer></script>`,
	Styles:  `<link rel="stylesheet" href="/⛰/Default.8912112e.css"><link rel="stylesheet" href="/⛰/Admin.db3e959a.css"><link rel="stylesheet" href="/⛰/User.68febbf8.css"><link rel="stylesheet" href="/⛰/User.09b09c46.css"><link rel="stylesheet" href="/⛰/Admin.9d294a13.css"><link rel="stylesheet" href="/⛰/AdminRebuild.38d394c2.css">`,	
}
