import { page } from '$app/state';
import { TestInstanceEventTypeEnum } from '$lib/api_types';
import { API } from '$lib/services/api.svelte';

// src/lib/ExamLogger.ts
export interface ExamLoggerOptions {
	endpoint?: string;
	sendIntervalMs?: number;
	maxBatchSize?: number;
	includeDebugWidget?: boolean;
	pageId?: string;
	instanceId?: string | null;
	debug?: boolean;
}

export interface ExamEvent {
	occuredAt: string; // ISO timestamp
	eventType: string; // event type
	eventData: any; // event detail payload
	pageId: string;
}

type Listener = [
	EventTarget,
	string,
	EventListenerOrEventListenerObject,
	boolean | AddEventListenerOptions | undefined
];

export default class ExamLogger {
	private config: Required<ExamLoggerOptions>;
	private buffer: ExamEvent[] = [];
	private lastSentAt = Date.now();
	private visStart: number | null;
	private debugEl: HTMLElement | null = null;
	private flushTimer: ReturnType<typeof setInterval> | null = null;
	private listeners: Listener[] = [];
	private isRunning: boolean = false;

	constructor(options: ExamLoggerOptions = {}) {
		this.config = Object.assign(
			{
				endpoint:
					import.meta.env.VITE_API_URL + `/api/v2/tests/${page.params.instanceId}/telemetry`,
				sendIntervalMs: 5000,
				maxBatchSize: 50,
				includeDebugWidget: false,
				pageId: (crypto as any).randomUUID
					? (crypto as any).randomUUID()
					: String(Date.now()) + Math.random(),
				instanceId: null,
				debug: false
			},
			options
		);

		this.visStart = document.visibilityState === 'visible' ? performance.now() : null;
	}

	public start() {
		this.isRunning = true;
		this.installListeners();
		this.flushTimer = setInterval(() => this.maybeFlush(), this.config.sendIntervalMs);

		this.record(TestInstanceEventTypeEnum.PAGELOAD, {
			url: location.pathname + location.search,
			w: window.innerWidth,
			h: window.innerHeight,
			tz: Intl.DateTimeFormat().resolvedOptions().timeZone
		});
	}

	public stop() {
		this.isRunning = false;
		this.removeListeners();
		if (this.flushTimer) clearInterval(this.flushTimer);
		this.maybeFlush(true);
		if (this.debugEl) this.debugEl.remove();
	}

	public running(): boolean {
		return this.isRunning;
	}

	public record(eventType: TestInstanceEventTypeEnum, eventData: Record<string, any> = {}) {
		const evt: ExamEvent = {
			occuredAt: new Date().toISOString(),
			eventType,
			eventData,
			pageId: this.config.pageId
		};
		this.buffer.push(evt);
		this.maybeFlush();
		if (this.config.debug) {
			this.debug(eventType, eventData);
		}
	}

	private maybeFlush(force = false) {
		console.log('tried to flush', this.buffer.length);
		if (this.buffer.length === 0) return;
		if (
			force ||
			this.buffer.length >= this.config.maxBatchSize ||
			Date.now() - this.lastSentAt >= this.config.sendIntervalMs
		) {
			const batch = this.buffer.splice(0, this.config.maxBatchSize);
			this.sendBatch(batch);
		}
	}

	private sendBatch(events: ExamEvent[]) {
		const ok = false;
		// navigator.sendBeacon &&
		// navigator.sendBeacon(this.config.endpoint, new Blob([payload], headers as any));

		if (!ok) {
			console.log(events);
			API.request(this.config.endpoint, {
				method: 'POST',
				body: {
					events: events
				},
				keepalive: true,
				credentials: 'include'
			}).catch((err) => {
				console.error('Failed to send telemetry', err, this.config.endpoint);
				this.record(TestInstanceEventTypeEnum.NETWORKERROR, {
					occuredAt: new Date().toISOString()
				});
				setTimeout(() => {
					this.sendBatch(events);
				}, 1500);
			});
		}
		this.lastSentAt = Date.now();
	}

	// --- Debug widget ---
	private debug(type: TestInstanceEventTypeEnum, detail: any) {
		if (!this.config.includeDebugWidget) return;
		if (!this.debugEl) {
			const el = document.createElement('div');
			el.id = 'exam-logger-debug';
			el.style.cssText =
				'position:fixed;right:12px;bottom:12px;z-index:99999;background:rgba(0,0,0,.75);color:#fff;font:12px system-ui;padding:10px;border-radius:10px;max-width:360px;max-height:200px;overflow:auto';
			el.innerHTML = `<strong>ExamLogger</strong><ul id="exam-logger-list"></ul>`;
			document.body.appendChild(el);
			this.debugEl = el;
		}
		const ul = this.debugEl.querySelector('#exam-logger-list') as HTMLElement;
		const li = document.createElement('li');
		li.textContent = `${new Date().toLocaleTimeString()} – ${type}`;
		ul.prepend(li);
		while (ul.children.length > 25) ul.removeChild(ul.lastChild!);
	}

	// --- Event listeners ---
	private installListeners() {
		const on = (
			target: EventTarget,
			type: string,
			fn: EventListenerOrEventListenerObject,
			opts?: boolean | AddEventListenerOptions
		) => {
			target.addEventListener(type, fn, opts);
			this.listeners.push([target, type, fn, opts]);
		};

		['copy', 'cut', 'paste'].forEach((type) => {
			on(
				document,
				type,
				async (evt: Event) => {
					const e = evt as ClipboardEvent;

					if (evt.type == 'paste') {
						const cb = e.clipboardData;
						const cbData: any = {};
						if (cb) {
							// Get all string types
							for (const type of cb.types) {
								if (type === 'Files') continue; // skip files here
								if (type === 'text/html') {
									cbData['text'] = cb.getData('text');
									continue;
								}
								cbData[type] = cb.getData(type);
							}

							cbData['text'] = cb.getData('text');

							// Get files with raw bytes
							const filesWithBytes: { file: File; bytes: string }[] = [];

							for (const item of cb.items) {
								if (item.kind === 'file') {
									const file = item.getAsFile();
									if (file) {
										const bytes = await file.arrayBuffer();
										filesWithBytes.push({
											file,
											bytes: btoa(String.fromCharCode(...new Uint8Array(bytes)))
										});
									}
								}
							}

							if (filesWithBytes.length > 0) {
								cbData['Files'] = filesWithBytes;
							}
						}

						this.record(TestInstanceEventTypeEnum.CLIPBOARD, {
							type: evt.type,
							clipboardData: cbData
						});
					} else {
						this.record(TestInstanceEventTypeEnum.CLIPBOARD, {
							type: evt.type
						});
					}
				},
				true
			);
		});

		on(document, 'contextmenu', () => this.record(TestInstanceEventTypeEnum.CONTEXTMENU, {}), true);
		on(document, 'selectstart', () => this.record(TestInstanceEventTypeEnum.SELECTSTART, {}), true);
		on(document, 'dragstart', () => this.record(TestInstanceEventTypeEnum.DRAGSTART, {}), true);
		on(document, 'drop', (e: Event) => this.record(TestInstanceEventTypeEnum.DROP, {}), true);

		const watched = new Set([
			'c',
			'x',
			'v',
			'a',
			'p',
			's',
			'l',
			'k',
			'f',
			'o',
			'n',
			't',
			'w',
			'r',
			'z',
			'y'
		]);
		on(
			document,
			'keydown',
			(evt: Event) => {
				const e = evt as KeyboardEvent;
				console.log(e.key);
				const key = e.key.toLowerCase();
				if ((e.ctrlKey || e.metaKey) && watched.has(key)) {
					this.record(TestInstanceEventTypeEnum.SHORTCUT, {
						key,
						ctrl: e.ctrlKey,
						meta: e.metaKey,
						alt: e.altKey,
						shift: e.shiftKey
					});
				}
				if (key === 'printscreen') this.record(TestInstanceEventTypeEnum.PRINTSCREEN, {});
			},
			true
		);

		on(window, 'beforeprint', () => this.record(TestInstanceEventTypeEnum.PRINT, {}));
		on(window, 'afterprint', () => this.record(TestInstanceEventTypeEnum.PRINT, {}));

		on(document, 'visibilitychange', () => {
			if (document.visibilityState === 'hidden') {
				const openMs =
					this.visStart == null ? null : Math.max(0, performance.now() - this.visStart);
				this.record(TestInstanceEventTypeEnum.TABHIDDEN, { openMs });
				this.visStart = null;
				this.maybeFlush(true);
			} else {
				this.visStart = performance.now();
				this.record(TestInstanceEventTypeEnum.TABVISIBLE, {});
			}
		});

		on(window, 'blur', () => this.record(TestInstanceEventTypeEnum.BLUR, {}));
		on(window, 'focus', () => this.record(TestInstanceEventTypeEnum.FOCUS, {}));
		on(document, 'fullscreenchange', () =>
			this.record(TestInstanceEventTypeEnum.FULLSCREEN, {
				isFullscreen: !!document.fullscreenElement
			})
		);

		let resizeTimer: any = null;
		on(window, 'resize', () => {
			clearTimeout(resizeTimer);
			resizeTimer = setTimeout(
				() =>
					this.record(TestInstanceEventTypeEnum.RESIZE, {
						w: window.innerWidth,
						h: window.innerHeight
					}),
				150
			);
		});

		on(window, TestInstanceEventTypeEnum.HIDE, () => this.maybeFlush(true));
		on(window, TestInstanceEventTypeEnum.UNLOAD, () => this.maybeFlush(true));
	}

	private removeListeners() {
		for (const [target, type, fn, opts] of this.listeners) {
			target.removeEventListener(type, fn, opts);
		}
		this.listeners = [];
	}
}
