CREATE TABLE `fund` (
    `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '序号',
    `code` varchar(32) NOT NULL DEFAULT '' COMMENT '基金代码',
    `name` varchar(64) NOT NULL DEFAULT '' COMMENT '基金简称',
    `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `unit_value` double DEFAULT 0.0 COMMENT '单位净值',
    `total_value` double DEFAULT 0.0 COMMENT '累计净值',
    `dayswell_rate` varchar(20) DEFAULT ''  COMMENT '日增长率',
    `weekswell_rate` varchar(20) DEFAULT ''  COMMENT '近一周增长率',
    `monthswell_rate` varchar(20) DEFAULT ''  COMMENT '近一个月增长率',
    `threemonthswell_rate` varchar(20) DEFAULT ''  COMMENT '近三个月增长率',
    `sixmonthswell_rate` varchar(20) DEFAULT ''  COMMENT '近六个月增长率',
    `yearswell_rate` varchar(20) DEFAULT ''  COMMENT '近一年增长率',
    `twoyearswell_rate` varchar(20) DEFAULT ''  COMMENT '近两年增长率',
    `threeyearswell_rate` varchar(20) DEFAULT ''  COMMENT '近三年增长率',
    `thisyearwell_rate` varchar(20) DEFAULT ''  COMMENT '近年来增长率',
    `createswell_rate` varchar(20) DEFAULT ''  COMMENT '成立以来',
    `custom_rate` varchar(200) DEFAULT ''  COMMENT '自定义',
    `fee` varchar(32) DEFAULT ''  COMMENT '手续费',
    `isbuy` int(3)  NOT NULL DEFAULT 0  COMMENT '是否可以购买',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT ='基金净值表';