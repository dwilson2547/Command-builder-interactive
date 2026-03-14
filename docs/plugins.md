# Plugin / Sharing System

Command Builder supports importing config files from any HTTP/HTTPS URL.  
This lets teams share command sets, and individuals publish their personal toolboxes.

---

## Importing a config

### From the search screen

```
/import https://raw.githubusercontent.com/example/repo/main/mytools.yaml
```

Press **Enter**. The config is fetched, validated, and saved to  
`~/.config/command-builder/configs/`.

### From the config manager screen

1. Open the config manager: type `/config` and press Enter.
2. Press **i** to open the import prompt.
3. Paste a URL and press **Enter**.

---

## Name collision handling

If the imported config's `name` field conflicts with an existing config, a numeric  
suffix is appended automatically (e.g. `my-tools-1`, `my-tools-2`, …).

---

## Hosting a shareable config

Any YAML file following the [config format](config-format.md) can be hosted at a  
publicly accessible URL and shared.

**Example — GitHub raw URL:**

```
https://raw.githubusercontent.com/yourname/command-builder-configs/main/k8s.yaml
```

**Example — self-hosted:**

```
https://tools.example.com/configs/devops.yaml
```

---

## Security considerations

- Only import configs from **trusted sources**.
- Configs define command *templates*, not executable code; the user always sees the  
  full command before it is printed.
- Imported configs are stored locally and can be reviewed or deleted at any time via  
  the config manager.
