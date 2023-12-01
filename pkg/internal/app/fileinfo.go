package app

type FileInfo struct {
	// ------------
	// Microsoft WOPI check file info specification:
	// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/checkfileinfo
	// ------------

	// The string name of the file, including extension, without a path. Used for display in user interface (UI), and determining the extension of the file.
	BaseFileName string `json:"BaseFileName,omitempty"`
	//A string that uniquely identifies the owner of the file. In most cases, the user who uploaded or created the file should be considered the owner.
	OwnerID string `json:"OwnerId,omitempty"`
	// A string value uniquely identifying the user currently accessing the file.
	UserID string `json:"UserId,omitempty"`
	// The size of the file in bytes, expressed as a long, a 64-bit signed integer.
	Size *int64 `json:"Size,omitempty"`
	// The current version of the file based on the server’s file version schema, as a string. This value must change when the file changes, and version values must never repeat for a given file.
	Version string `json:"Version,omitempty"`
	// A string that is the name of the user, suitable for displaying in UI.
	UserFriendlyName string `json:"UserFriendlyName,omitempty"`
	// A 256 bit SHA-2-encoded [FIPS 180-2] hash of the file contents, as a Base64-encoded string. Used for caching purposes in WOPI clients.
	SHA256 string `json:"SHA256,omitempty"`
	// A string value representing the file extension for the file. This value must begin with a .. If provided, WOPI clients will use this value as the file extension. Otherwise the extension will be parsed from the BaseFileName.
	FileExtension string `json:"FileExtension,omitempty"`
	// An integer value that indicates the maximum length for file names that the WOPI host supports, excluding the file extension. The default value is 250. Note that WOPI clients will use this default value if the property is omitted or if it is explicitly set to 0.
	FileNameMaxLength int `json:"FileNameMaxLength,omitempty"`
	// A string that represents the last time that the file was modified. This time must always be a must be a UTC time, and must be formatted in ISO 8601 round-trip format. For example, "2009-06-15T13:45:30.0000000Z".
	LastModifiedTime string `json:"LastModifiedTime,omitempty"`
	// A string value containing information about the user. This string can be passed from a WOPI client to the host by means of a PutUserInfo operation. If the host has a UserInfo string for the user, they must include it in this property. See the PutUserInfo documentation for more details.
	UserInfo string `json:"UserInfo,omitempty"`
	// A Boolean value that indicates that, for this user, the file cannot be changed.
	ReadOnly bool `json:"ReadOnly"`
	// A Boolean value that indicates that if host is temporarily unable to process writes on a file
	TemporarilyNotWritable bool `json:"TemporarilyNotWritable,omitempty"`
	// A string value indicating whether the current document is shared with other users. The value can change upon adding or removing permissions to other users. Clients should use this value to help decide when to enable collaboration features as a document must be Shared in order to multi-user collaboration on the document.
	SharingStatus string `json:"SharingStatus,omitempty"`
	// A Boolean value that indicates that the WOPI client should restrict what actions the user can perform on the file. The behavior of this property is dependent on the WOPI client.
	RestrictedWebViewOnly bool `json:"RestrictedWebViewOnly"`
	// A Boolean value that indicates that the user has permission to view a broadcast of this file.
	UserCanAttend bool `json:"UserCanAttend"`
	// A Boolean value that indicates that the user has permission to broadcast this file to a set of users who have permission to broadcast or view a broadcast of the current file.
	UserCanPresent bool `json:"UserCanPresent"`
	// A Boolean value that indicates the user does not have sufficient permission to create new files on the WOPI server. Setting this to true tells the WOPI client that calls to PutRelativeFile will fail for this user on the current file.
	UserCanNotWriteRelative bool `json:"UserCanNotWriteRelative"`
	// A Boolean value that indicates the user has permission to rename the current file.
	UserCanRename bool `json:"UserCanRename"`
	// A Boolean value that indicates that the user has permission to alter the file. Setting this to true tells the WOPI client that it can call PutFile on behalf of the user.
	UserCanWrite bool `json:"UserCanWrite"`
	// A Boolean value that indicates the WOPI client should close the window or tab when the user activates any Close UI in the WOPI client.
	CloseButtonClosesWindow bool `json:"CloseButtonClosesWindow"`
	// A string value indicating whether the WOPI client should disable Copy and Paste functionality within the application. The default is to permit all Copy and Paste functionality, i.e. the setting has no effect.
	CopyPasteRestrictions string `json:"CopyPasteRestrictions,omitempty"`
	// A Boolean value that indicates the WOPI client should disable all print functionality.
	DisablePrint bool `json:"DisablePrint"`
	// A Boolean value that indicates the WOPI client should disable all machine translation functionality.
	DisableTranslation bool `json:"DisableTranslation"`
	// A URI to a web page that the WOPI client should navigate to when the application closes, or in the event of an unrecoverable error.
	CloseUrl string `json:"CloseUrl,omitempty"`
	// A user-accessible URI to the file intended to allow the user to download a copy of the file.
	DownloadUrl string `json:"DownloadUrl,omitempty"`
	// A URI to a location that allows the user to create an embeddable URI to the file.
	FileEmbedCommandUrl string `json:"FileEmbedCommandUrl,omitempty"`
	// A URI to a location that allows the user to share the file.
	FileSharingUrl string `json:"FileSharingUrl,omitempty"`
	// A URI to the file location that the WOPI client uses to get the file. If this is provided, the WOPI client may use this URI to get the file instead of a GetFile request. A host might set this property if it is easier or provides better performance to serve files from a different domain than the one handling standard WOPI requests. WOPI clients must not add or remove parameters from the URL; no other parameters, including the access token, should be appended to the FileUrl before it is used.
	FileUrl string `json:"FileUrl,omitempty"`
	// A URI to a location that allows the user to view the version history for the file.
	FileVersionUrl string `json:"FileVersionUrl,omitempty"`
	// A URI to a host page that loads the edit WOPI action.
	HostEditUrl string `json:"HostEditUrl,omitempty"`
	// A URI to a web page that provides access to a viewing experience for the file that can be embedded in another HTML page. This is typically a URI to a host page that loads the embedview WOPI action.
	HostEmbeddedViewUrl string `json:"HostEmbeddedViewUrl,omitempty"`
	// A URI to a host page that loads the view WOPI action. This URL is used by Office Online to navigate between view and edit mode.
	HostViewUrl string `json:"HostViewUrl,omitempty"`
	// A URI that will sign the current user out of the host’s authentication system.
	SignoutUrl string `json:"SignoutUrl,omitempty"`
	// A string that indicates the brand name of the host.
	BreadcrumbBrandName string `json:"BreadcrumbBrandName,omitempty"`
	// A URI to a web page that the WOPI client should navigate to when the user clicks on UI that displays BreadcrumbBrandName.
	BreadcrumbBrandUrl string `json:"BreadcrumbBrandUrl,omitempty"`
	// A string that indicates the name of the file. If this is not provided, WOPI clients may use the BaseFileName value.
	BreadcrumbDocName string `json:"BreadcrumbDocName,omitempty"`
	// A string that indicates the name of the container that contains the file.
	BreadcrumbFolderName string `json:"BreadcrumbFolderName,omitempty"`
	// A URI to a web page that the WOPI client should navigate to when the user clicks on UI that displays BreadcrumbFolderName.
	BreadcrumbFolderUrl string `json:"BreadcrumbFolderUrl,omitempty"`
	// A Boolean value that indicates a WOPI client may connect to Microsoft services to provide end-user functionality.
	AllowAdditionalMicrosoftServices bool `json:"AllowAdditionalMicrosoftServices"`
	// A Boolean value that indicates that in the event of an error, the WOPI client is permitted to prompt the user for permission to collect a detailed report about their specific error. The information gathered could include the user’s file and other session-specific state.
	AllowErrorReportPrompt bool `json:"AllowErrorReportPrompt"`
	// A Boolean value that indicates a WOPI client may allow connections to external services referenced in the file (for example, a marketplace of embeddable JavaScript apps).
	AllowExternalMarketplace bool `json:"AllowExternalMarketplace"`
	// A string value offering guidance to the WOPI client as to how to differentiate client throttling behaviors between the user and documents combinations from the WOPI host.
	ClientThrottlingProtection string `json:"ClientThrottlingProtection,omitempty"`
	// A string value indicating whether the WOPI host is experiencing capacity problems and would like to reduce the frequency at which the WOPI clients make calls to the host
	RequestedCallThrottling string `json:"RequestedCallThrottling,omitempty"`
	// An array of strings containing the Share URL types supported by the host.
	SupportedShareUrlTypes []string `json:"SupportedShareUrlTypes,omitempty"`
	// A Boolean value that indicates that the host supports the following WOPI operations: ExecuteCellStorageRequest, ExecuteCellStorageRelativeRequest
	SupportsCobalt bool `json:"SupportsCobalt"`
	// A Boolean value that indicates that the host supports the following WOPI operations: CheckContainerInfo, CreateChildContainer, CreateChildFile, DeleteContainer, DeleteFile, EnumerateAncestors (containers), EnumerateAncestors (files), EnumerateChildren (containers), GetEcosystem (containers), RenameContainer
	SupportsContainers bool `json:"SupportsContainers"`
	// A Boolean value that indicates that the host supports the following WOPI operations: CheckEcosystem, GetEcosystem (containers), GetEcosystem (files), GetRootContainer (ecosystem)
	SupportsEcosystem bool `json:"SupportsEcosystem"`
	// A Boolean value that indicates that the host supports lock IDs up to 1024 ASCII characters long. If not provided, WOPI clients will assume that lock IDs are limited to 256 ASCII characters.
	SupportsExtendedLockLength bool `json:"SupportsExtendedLockLength"`
	// A Boolean value that indicates that the host supports the following WOPI operations: CheckFolderInfo, EnumerateChildren (folders), DeleteFile
	SupportsFolders bool `json:"SupportsFolders"`
	// A Boolean value that indicates that the host supports the GetFileWopiSrc (ecosystem) operation.
	SupportsGetFileWopiSrc bool `json:"SupportsGetFileWopiSrc"`
	// A Boolean value that indicates that the host supports the GetLock operation.
	SupportsGetLock bool `json:"SupportsGetLock"`
	// A Boolean value that indicates that the host supports the following WOPI operations: Lock, Unlock, RefreshLock, UnlockAndRelock operations for this file.
	SupportsLocks bool `json:"SupportsLocks"`
	// A Boolean value that indicates that the host supports the RenameFile operation.
	SupportsRename bool `json:"SupportsRename"`
	// A Boolean value that indicates that the host supports the following WOPI operations: PutFile, PutRelativeFile
	SupportsUpdate bool `json:"SupportsUpdate"` // whether "Putfile" and "PutRelativeFile" work
	// A Boolean value that indicates that the host supports the DeleteFile operation.
	SupportsDeleteFile bool `json:"SupportsDeleteFile"`
	// A Boolean value that indicates that the host supports the PutUserInfo operation.
	SupportsUserInfo bool `json:"SupportsUserInfo"`
	// A Boolean value indicating whether the user is authenticated with the host or not. Hosts should always set this to true for unauthenticated users, so that clients are aware that the user is anonymous. When setting this to true, hosts can choose to omit the UserId property, but must still set the OwnerId property.
	IsAnonymousUser bool `json:"IsAnonymousUser,omitempty"`
	// A Boolean value indicating whether the user is an education user or not.
	IsEduUser bool `json:"IsEduUser,omitempty"`
	// A Boolean value indicating whether the user is a business user or not.
	LicenseCheckForEditIsEnabled bool `json:"LicenseCheckForEditIsEnabled,omitempty"`
	// A Boolean value indicating whether activities are supported
	SupportsAddActivities bool `json:"SupportsAddActivities,omitempty"`

	// ------------
	// Collabora WOPI check file info specification:
	// https://sdk.collaboraonline.com/docs/advanced_integration.html
	// ------------

	// If set to true, it allows the document owner (the one with OwnerId =UserId) to send a closedocument message (see protocol.txt)
	EnableOwnerTermination bool `json:"EnableOwnerTermination,omitempty"`
	// Disables export functionality in backend. If set to true, HideExportOption is assumed to be true
	DisableExport bool `json:"DisableExport,omitempty"`
	// Disables copying from the document in libreoffice online backend. Pasting into the document would still be possible. However, it is still possible to do an “internal” cut/copy/paste.
	DisableCopy bool `json:"DisableCopy,omitempty"`
}
