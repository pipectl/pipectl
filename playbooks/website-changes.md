# Playbook: Website Changes

## Goal

Change or extend the `website/` VitePress site without reintroducing dependence on
network-fetched, third-party JavaScript.

## Why This Constraint Exists

Some visitors are behind corporate proxies (e.g. Netskope) that do TLS interception and
rewrite/intercept JavaScript in transit. This makes JS-heavy sites — and especially sites
that pull JS from third-party CDNs — noticeably slow or broken for those visitors. The
site currently avoids this: no analytics/tracking scripts, no CDN-hosted search widget,
self-hosted fonts, and all page content is prerendered to static HTML so it's readable
even if JS execution is degraded.

## Rules

- Do not add analytics or tracking scripts (Google Analytics, Plausible, Fathom, etc.).
- Do not switch search from VitePress's local provider (`search: { provider: 'local' }`
  in `website/.vitepress/config.mts`) to Algolia DocSearch or any other CDN-hosted search.
- Do not load fonts from `fonts.googleapis.com` / `fonts.gstatic.com` — self-host font
  files instead, as the default VitePress theme already does.
- Do not embed externally-hosted widgets (chat widgets, third-party comment systems, ad
  scripts, etc.).
- Prefer content that renders as static HTML over content that depends on client-side JS
  to appear at all.

## Enforcement

`.github/workflows/docs.yml` greps `website/docs` and `website/.vitepress` for known
offending patterns (Google Fonts domains, DocSearch/Algolia, common analytics domains)
before building, and fails the build if any are found. Update that pattern list if a new
category of banned third-party JS is identified.
